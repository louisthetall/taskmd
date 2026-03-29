import { NavLink } from "react-router-dom";
import { useConfig } from "../../hooks/use-config.ts";
import { useProject } from "../../hooks/use-project.ts";

interface Tab {
  path: string;
  label: string;
}

const baseTabs: Tab[] = [
  { path: "/tasks", label: "Tasks" },
  { path: "/next", label: "Next Up" },
  { path: "/board", label: "Board" },
  { path: "/tracks", label: "Tracks" },
  { path: "/graph", label: "Graph" },
  { path: "/feed", label: "Activity" },
  { path: "/stats", label: "Stats" },
  { path: "/validate", label: "Validate" },
];

function useTabs(): Tab[] {
  const { project } = useProject();
  const { phases } = useConfig(project);
  if (phases.length === 0) return baseTabs;
  return [
    ...baseTabs.slice(0, 5), // Tasks through Graph
    { path: "/phases", label: "Phases" },
    ...baseTabs.slice(5),    // Stats, Validate
  ];
}

const navLinkClass = ({ isActive }: { isActive: boolean }) =>
  `px-3 py-1.5 text-sm rounded-md transition-colors ${
    isActive
      ? "bg-gray-900 text-white dark:bg-white dark:text-gray-900"
      : "text-gray-600 hover:text-gray-900 hover:bg-gray-100 dark:text-gray-400 dark:hover:text-gray-100 dark:hover:bg-gray-700"
  }`;

export function DesktopNav() {
  const tabs = useTabs();
  return (
    <nav className="flex items-center gap-1" data-arrow-nav>
      {tabs.map((tab) => (
        <NavLink key={tab.path} to={tab.path} className={navLinkClass}>
          {tab.label}
        </NavLink>
      ))}
    </nav>
  );
}

interface HeaderActionsProps {
  onSearchOpen: () => void;
}

export function HeaderActions({ onSearchOpen }: HeaderActionsProps) {
  return (
    <div className="hidden md:flex items-center gap-1">
      <button
        onClick={onSearchOpen}
        className="px-2.5 py-1.5 text-sm rounded-md transition-colors text-gray-600 hover:text-gray-900 hover:bg-gray-100 dark:text-gray-400 dark:hover:text-gray-100 dark:hover:bg-gray-700 flex items-center gap-1.5"
        aria-label="Search tasks"
      >
        <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <kbd className="text-[10px] font-medium text-gray-400 bg-gray-100 dark:bg-gray-700 dark:text-gray-500 px-1 py-0.5 rounded">
          ⌘K
        </kbd>
      </button>
      <a
        href="https://driangle.github.io/taskmd/"
        target="_blank"
        rel="noopener noreferrer"
        className="px-2 py-1.5 text-sm rounded-md transition-colors text-gray-600 hover:text-gray-900 hover:bg-gray-100 dark:text-gray-400 dark:hover:text-gray-100 dark:hover:bg-gray-700"
      >
        Docs ↗
      </a>
      <a
        href="https://github.com/driangle/taskmd"
        target="_blank"
        rel="noopener noreferrer"
        className="p-1.5 rounded-md text-gray-600 hover:text-gray-900 hover:bg-gray-100 dark:text-gray-400 dark:hover:text-gray-100 dark:hover:bg-gray-700 transition-colors"
        aria-label="GitHub repository"
      >
        <svg className="w-4 h-4" viewBox="0 0 16 16" fill="currentColor">
          <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z" />
        </svg>
      </a>
    </div>
  );
}

export function MobileMenu() {
  const tabs = useTabs();
  return (
    <nav className="md:hidden border-t border-gray-200 dark:border-gray-700 px-4 py-2 space-y-1">
      {tabs.map((tab) => (
        <NavLink
          key={tab.path}
          to={tab.path}
          className={({ isActive }) =>
            `block px-3 py-2 text-sm rounded-md transition-colors ${
              isActive
                ? "bg-gray-900 text-white dark:bg-white dark:text-gray-900"
                : "text-gray-600 hover:text-gray-900 hover:bg-gray-100 dark:text-gray-400 dark:hover:text-gray-100 dark:hover:bg-gray-700"
            }`
          }
        >
          {tab.label}
        </NavLink>
      ))}
      <a
        href="https://driangle.github.io/taskmd/"
        target="_blank"
        rel="noopener noreferrer"
        className="block px-3 py-2 text-sm rounded-md transition-colors text-gray-600 hover:text-gray-900 hover:bg-gray-100 dark:text-gray-400 dark:hover:text-gray-100 dark:hover:bg-gray-700"
      >
        Docs ↗
      </a>
      <a
        href="https://github.com/driangle/taskmd"
        target="_blank"
        rel="noopener noreferrer"
        className="block px-3 py-2 text-sm rounded-md transition-colors text-gray-600 hover:text-gray-900 hover:bg-gray-100 dark:text-gray-400 dark:hover:text-gray-100 dark:hover:bg-gray-700"
      >
        GitHub ↗
      </a>
    </nav>
  );
}
