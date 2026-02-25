// Minimal vscode mock for tests that import modules depending on vscode
export class Range {
  constructor(
    public startLine: number,
    public startCharacter: number,
    public endLine: number,
    public endCharacter: number
  ) {}
}

export class Position {
  constructor(public line: number, public character: number) {}
}

export enum DiagnosticSeverity {
  Error = 0,
  Warning = 1,
  Information = 2,
  Hint = 3,
}

export class Diagnostic {
  source?: string;
  constructor(
    public range: Range,
    public message: string,
    public severity?: DiagnosticSeverity
  ) {}
}

export enum CompletionItemKind {
  EnumMember = 19,
  Value = 11,
  Reference = 17,
}

export class CompletionItem {
  insertText?: string;
  range?: Range;
  detail?: string;
  constructor(public label: string, public kind?: CompletionItemKind) {}
}
