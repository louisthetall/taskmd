import { useState, useRef, useEffect, useMemo } from "react";
import { toggleInSet } from "../tasks/TaskTable/utils.ts";

interface TagAutocompleteProps {
  availableTags: string[];
  selectedTags: Set<string>;
  onTagsChange: (next: Set<string>) => void;
}

export function TagAutocomplete({
  availableTags,
  selectedTags,
  onTagsChange,
}: TagAutocompleteProps) {
  const [query, setQuery] = useState("");
  const [open, setOpen] = useState(false);
  const [activeIndex, setActiveIndex] = useState(0);
  const containerRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const listRef = useRef<HTMLUListElement>(null);

  const suggestions = useMemo(() => {
    const unselected = availableTags.filter((t) => !selectedTags.has(t));
    if (!query) return unselected;
    const lower = query.toLowerCase();
    return unselected.filter((t) => t.toLowerCase().includes(lower));
  }, [availableTags, selectedTags, query]);

  // Reset active index when suggestions change
  useEffect(() => {
    setActiveIndex(0);
  }, [suggestions]);

  // Scroll active item into view
  useEffect(() => {
    if (!open || !listRef.current) return;
    const item = listRef.current.children[activeIndex] as HTMLElement | undefined;
    item?.scrollIntoView({ block: "nearest" });
  }, [activeIndex, open]);

  // Close on outside click
  useEffect(() => {
    function handleClick(e: MouseEvent) {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClick);
    return () => document.removeEventListener("mousedown", handleClick);
  }, []);

  function selectTag(tag: string) {
    onTagsChange(toggleInSet(selectedTags, tag));
    setQuery("");
    setOpen(false);
    inputRef.current?.focus();
  }

  function removeTag(tag: string) {
    onTagsChange(toggleInSet(selectedTags, tag));
  }

  function handleKeyDown(e: React.KeyboardEvent) {
    if (!open) {
      if (e.key === "ArrowDown" || e.key === "Enter") {
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
          selectTag(suggestions[activeIndex]);
        }
        break;
      case "Escape":
        e.preventDefault();
        setOpen(false);
        break;
    }
  }

  return (
    <div className="flex items-start gap-2 flex-wrap">
      <span className="text-xs text-gray-500 dark:text-gray-400 font-medium mt-1.5">
        Tags:
      </span>

      {/* Selected tag chips */}
      {[...selectedTags].map((tag) => (
        <span
          key={tag}
          className="inline-flex items-center gap-1 px-2.5 py-1 text-xs rounded-full bg-blue-100 text-blue-700 ring-1 ring-blue-300 dark:bg-blue-900/30 dark:text-blue-300 dark:ring-blue-700"
        >
          {tag}
          <button
            onClick={() => removeTag(tag)}
            className="hover:text-blue-900 dark:hover:text-blue-100"
            aria-label={`Remove ${tag} filter`}
          >
            &times;
          </button>
        </span>
      ))}

      {/* Autocomplete input */}
      <div ref={containerRef} className="relative">
        <input
          ref={inputRef}
          type="text"
          value={query}
          onChange={(e) => {
            setQuery(e.target.value);
            setOpen(true);
          }}
          onFocus={() => setOpen(true)}
          onKeyDown={handleKeyDown}
          placeholder="Add tag..."
          className="px-2.5 py-1 text-xs border border-gray-200 rounded-full bg-gray-50 text-gray-700 placeholder-gray-400 focus:outline-none focus:ring-1 focus:ring-blue-300 focus:border-blue-300 dark:bg-gray-800/50 dark:border-gray-700 dark:text-gray-300 dark:placeholder-gray-500 dark:focus:ring-blue-700 dark:focus:border-blue-700 w-36"
          role="combobox"
          aria-expanded={open}
          aria-autocomplete="list"
          aria-activedescendant={open && suggestions[activeIndex] ? `tag-option-${activeIndex}` : undefined}
        />

        {open && suggestions.length > 0 && (
          <ul
            ref={listRef}
            role="listbox"
            className="absolute z-50 mt-1 left-0 w-48 max-h-48 overflow-auto rounded-md border border-gray-200 bg-white shadow-lg dark:bg-gray-800 dark:border-gray-700"
          >
            {suggestions.map((tag, i) => (
              <li
                key={tag}
                id={`tag-option-${i}`}
                role="option"
                aria-selected={i === activeIndex}
                onMouseDown={(e) => e.preventDefault()}
                onClick={() => selectTag(tag)}
                onMouseEnter={() => setActiveIndex(i)}
                className={`px-3 py-1.5 text-xs cursor-pointer ${
                  i === activeIndex
                    ? "bg-blue-50 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300"
                    : "text-gray-700 dark:text-gray-300"
                }`}
              >
                {tag}
              </li>
            ))}
          </ul>
        )}

        {open && query && suggestions.length === 0 && (
          <div className="absolute z-50 mt-1 left-0 w-48 rounded-md border border-gray-200 bg-white shadow-lg dark:bg-gray-800 dark:border-gray-700 px-3 py-2 text-xs text-gray-400">
            No matching tags
          </div>
        )}
      </div>
    </div>
  );
}
