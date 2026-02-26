import { useState, useRef, useEffect, useMemo } from "react";

interface FolderAutocompleteProps {
  folders: string[];
  value: string;
  onChange: (value: string) => void;
}

const MAX_SUGGESTIONS = 5;

/** Score a folder against a query. Lower is better. Returns -1 for no match. */
function similarityScore(folder: string, query: string): number {
  const f = folder.toLowerCase();
  const q = query.toLowerCase();
  // Exact match
  if (f === q) return 0;
  // Starts with query
  if (f.startsWith(q)) return 1;
  // Contains query
  const idx = f.indexOf(q);
  if (idx >= 0) return 2 + idx;
  return -1;
}

export function FolderAutocomplete({
  folders,
  value,
  onChange,
}: FolderAutocompleteProps) {
  const [query, setQuery] = useState(value);
  const [open, setOpen] = useState(false);
  const [activeIndex, setActiveIndex] = useState(0);
  const containerRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const listRef = useRef<HTMLUListElement>(null);

  // Sync external value changes into the input
  useEffect(() => {
    setQuery(value);
  }, [value]);

  const suggestions = useMemo(() => {
    if (!query) return folders.slice(0, MAX_SUGGESTIONS);
    return folders
      .map((f) => ({ folder: f, score: similarityScore(f, query) }))
      .filter((r) => r.score >= 0)
      .sort((a, b) => a.score - b.score)
      .slice(0, MAX_SUGGESTIONS)
      .map((r) => r.folder);
  }, [folders, query]);

  // eslint-disable-next-line react-hooks/set-state-in-effect
  useEffect(() => {
    setActiveIndex(0);
  }, [suggestions]);

  useEffect(() => {
    if (!open || !listRef.current) return;
    const item = listRef.current.children[activeIndex] as
      | HTMLElement
      | undefined;
    item?.scrollIntoView?.({ block: "nearest" });
  }, [activeIndex, open]);

  // Close on outside click
  useEffect(() => {
    function handleClick(e: MouseEvent) {
      if (
        containerRef.current &&
        !containerRef.current.contains(e.target as Node)
      ) {
        setOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClick);
    return () => document.removeEventListener("mousedown", handleClick);
  }, []);

  function selectFolder(folder: string) {
    setQuery(folder);
    onChange(folder);
    setOpen(false);
  }

  function handleKeyDown(e: React.KeyboardEvent) {
    if (!open) {
      if (e.key === "ArrowDown") {
        setOpen(true);
        e.preventDefault();
      }
      return;
    }

    switch (e.key) {
      case "ArrowDown":
        e.preventDefault();
        setActiveIndex((i) => Math.min(i + 1, suggestions.length - 1));
        break;
      case "ArrowUp":
        e.preventDefault();
        setActiveIndex((i) => Math.max(i - 1, 0));
        break;
      case "Enter":
        e.preventDefault();
        if (suggestions[activeIndex]) {
          selectFolder(suggestions[activeIndex]);
        }
        break;
      case "Escape":
        e.preventDefault();
        setOpen(false);
        break;
    }
  }

  function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
    const v = e.target.value;
    setQuery(v);
    setOpen(true);
    // Clear the filter when the input is emptied
    if (!v) onChange("");
  }

  return (
    <div ref={containerRef} className="relative">
      <input
        ref={inputRef}
        id="group-filter"
        type="text"
        value={query}
        onChange={handleChange}
        onFocus={() => setOpen(true)}
        onKeyDown={handleKeyDown}
        autoComplete="off"
        placeholder="All folders"
        className="min-h-[44px] sm:min-h-0 text-sm rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 px-2 py-1 w-40 placeholder-gray-400 dark:placeholder-gray-500 focus:outline-none focus:ring-1 focus:ring-blue-300 focus:border-blue-300 dark:focus:ring-blue-700 dark:focus:border-blue-700"
        role="combobox"
        aria-expanded={open}
        aria-autocomplete="list"
        aria-activedescendant={
          open && suggestions[activeIndex]
            ? `folder-option-${activeIndex}`
            : undefined
        }
      />

      {value && (
        <button
          type="button"
          onClick={() => {
            setQuery("");
            onChange("");
            inputRef.current?.focus();
          }}
          className="absolute right-1.5 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 text-sm leading-none"
          aria-label="Clear folder filter"
        >
          &times;
        </button>
      )}

      {open && suggestions.length > 0 && (
        <ul
          ref={listRef}
          role="listbox"
          className="absolute z-50 mt-1 left-0 w-48 max-h-48 overflow-auto rounded-md border border-gray-200 bg-white shadow-lg dark:bg-gray-800 dark:border-gray-700"
        >
          {suggestions.map((folder, i) => (
            <li
              key={folder}
              id={`folder-option-${i}`}
              role="option"
              aria-selected={i === activeIndex}
              onMouseDown={(e) => e.preventDefault()}
              onClick={() => selectFolder(folder)}
              onMouseEnter={() => setActiveIndex(i)}
              className={`px-3 py-1.5 text-sm cursor-pointer ${
                i === activeIndex
                  ? "bg-blue-50 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300"
                  : "text-gray-700 dark:text-gray-300"
              }`}
            >
              {folder}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
