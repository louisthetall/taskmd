import * as vscode from "vscode";
import { parseFrontmatter } from "./frontmatter";
import { validate } from "./validate";
import { toDiagnostics } from "./diagnostics";
import { TaskmdCompletionProvider } from "./completion";
import { isUnderTaskDir } from "./config";

const DIAGNOSTIC_SOURCE = "taskmd";

let diagnosticCollection: vscode.DiagnosticCollection;

export function activate(context: vscode.ExtensionContext): void {
  diagnosticCollection = vscode.languages.createDiagnosticCollection(DIAGNOSTIC_SOURCE);
  context.subscriptions.push(diagnosticCollection);

  // Validate on open
  context.subscriptions.push(
    vscode.workspace.onDidOpenTextDocument((doc) => validateDocument(doc))
  );

  // Validate on change
  context.subscriptions.push(
    vscode.workspace.onDidChangeTextDocument((e) => validateDocument(e.document))
  );

  // Clear on close
  context.subscriptions.push(
    vscode.workspace.onDidCloseTextDocument((doc) => {
      diagnosticCollection.delete(doc.uri);
    })
  );

  // Register completion provider
  const completionProvider = new TaskmdCompletionProvider();
  context.subscriptions.push(
    vscode.languages.registerCompletionItemProvider(
      { language: "markdown", scheme: "file" },
      completionProvider,
      ":", "-", "[", '"'
    )
  );

  // Watch for task file changes to invalidate the task ID cache
  const watcher = vscode.workspace.createFileSystemWatcher("**/*.md");
  watcher.onDidCreate(() => completionProvider.invalidateTaskCache());
  watcher.onDidDelete(() => completionProvider.invalidateTaskCache());
  watcher.onDidChange(() => completionProvider.invalidateTaskCache());
  context.subscriptions.push(watcher);

  // Validate already-open documents
  for (const doc of vscode.workspace.textDocuments) {
    validateDocument(doc);
  }
}

export function deactivate(): void {
  diagnosticCollection?.dispose();
}

function isTaskFile(doc: vscode.TextDocument): boolean {
  if (doc.languageId !== "markdown") return false;
  return isUnderTaskDir(doc.uri.fsPath);
}

function validateDocument(doc: vscode.TextDocument): void {
  if (!isTaskFile(doc)) {
    diagnosticCollection.delete(doc.uri);
    return;
  }

  const text = doc.getText();
  const parsed = parseFrontmatter(text);

  if (!parsed) {
    diagnosticCollection.set(doc.uri, []);
    return;
  }

  // Convert YAML parse errors to diagnostics
  const yamlDiags = parsed.errors.map((err) => {
    const diag = new vscode.Diagnostic(
      new vscode.Range(err.range.start.line, err.range.start.col, err.range.end.line, err.range.end.col),
      err.message,
      vscode.DiagnosticSeverity.Error
    );
    diag.source = DIAGNOSTIC_SOURCE;
    return diag;
  });

  const issues = validate(parsed);
  const issueDiags = toDiagnostics(issues, parsed);

  diagnosticCollection.set(doc.uri, [...yamlDiags, ...issueDiags]);
}
