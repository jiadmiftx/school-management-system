<!--
  SearchBar Component
  ===================
  Reusable search bar with input, button, and stats display.
  
  Usage:
  <SearchBar
    placeholder="Cari nama atau NIP..."
    bind:value={searchQuery}
    totalItems={totalItems}
    itemLabel="guru"
    on:search={handleSearch}
  />
-->
<script lang="ts">
    import { createEventDispatcher } from "svelte";

    export let placeholder: string = "Cari...";
    export let value: string = "";
    export let totalItems: number = 0;
    export let itemLabel: string = "item";
    export let showStats: boolean = true;

    const dispatch = createEventDispatcher();

    function handleSearch() {
        dispatch("search", { value });
    }

    function handleKeypress(e: KeyboardEvent) {
        if (e.key === "Enter") {
            handleSearch();
        }
    }
</script>

<div class="search-bar glass-card">
    <div class="search-input-wrapper">
        <span class="search-icon">
            <svg
                xmlns="http://www.w3.org/2000/svg"
                width="18"
                height="18"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
            >
                <circle cx="11" cy="11" r="8" />
                <path d="m21 21-4.35-4.35" />
            </svg>
        </span>
        <input
            type="text"
            {placeholder}
            bind:value
            on:keypress={handleKeypress}
        />
    </div>
    <button class="btn-search" on:click={handleSearch}>Cari</button>
    {#if showStats}
        <div class="stats">
            <span class="stat-item">
                <strong>{totalItems}</strong>
                {itemLabel} terdaftar
            </span>
        </div>
    {/if}
</div>

<style>
    .search-bar {
        display: flex;
        align-items: center;
        gap: 1rem;
        padding: 1rem 1.5rem;
        margin-bottom: 1.5rem;
        background: #ffffff;
        border: 1px solid #e5e7eb;
        border-radius: 16px;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
    }

    .search-input-wrapper {
        flex: 1;
        display: flex;
        align-items: center;
        background: var(--color-bg-input, #ffffff);
        border: 1px solid var(--color-border, #e5e7eb);
        border-radius: 10px;
        padding: 0 1rem;
    }

    .search-icon {
        font-size: 1.1rem;
        margin-right: 0.5rem;
        color: var(--color-text-muted, #9ca3af);
        display: flex;
    }

    .search-input-wrapper input {
        flex: 1;
        background: transparent;
        border: none;
        padding: 0.75rem 0;
        color: var(--color-text-primary, #111827);
        font-size: 1rem;
        outline: none;
    }

    .search-input-wrapper input::placeholder {
        color: var(--color-text-muted, #9ca3af);
    }

    .btn-search {
        background: var(--color-primary, #00ced1);
        border: none;
        color: white;
        padding: 0.75rem 1.5rem;
        border-radius: 10px;
        cursor: pointer;
        font-weight: 500;
        transition: background 0.2s;
    }

    .btn-search:hover {
        background: var(--color-primary-hover, #00b5b8);
    }

    .stats {
        color: var(--color-text-secondary, #6b7280);
        font-size: 0.875rem;
    }

    .stat-item strong {
        color: var(--color-text-primary, #111827);
    }

    @media (max-width: 768px) {
        .search-bar {
            flex-direction: column;
            align-items: stretch;
        }

        .stats {
            text-align: center;
        }
    }
</style>
