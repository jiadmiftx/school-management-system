<!--
  ActionButton Component
  ======================
  Reusable action button for tables (view, edit, delete, reset).
  
  Usage:
  <ActionButton type="view" href="/path/to/detail" title="Lihat Detail" />
  <ActionButton type="edit" on:click={() => openEditModal(item)} title="Edit" />
  <ActionButton type="delete" on:click={() => confirmDelete(item)} title="Hapus" />
-->
<script lang="ts">
    import { createEventDispatcher } from "svelte";

    export let type: "view" | "edit" | "delete" | "reset" = "view";
    export let href: string = "";
    export let title: string = "";

    const dispatch = createEventDispatcher();

    function handleClick(e: MouseEvent) {
        if (!href) {
            e.preventDefault();
            dispatch("click");
        }
    }

    // Icon SVG paths
    const icons = {
        view: `<circle cx="12" cy="12" r="3"/><path d="M12 5c-7 0-10 7-10 7s3 7 10 7 10-7 10-7-3-7-10-7z"/>`,
        edit: `<path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>`,
        delete: `<path d="M3 6h18"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>`,
        reset: `<path d="M3 12a9 9 0 0 1 9-9 9.75 9.75 0 0 1 6.74 2.74L21 8"/><path d="M21 3v5h-5"/><path d="M21 12a9 9 0 0 1-9 9 9.75 9.75 0 0 1-6.74-2.74L3 16"/><path d="M8 16H3v5"/>`,
    };
</script>

{#if href}
    <a {href} class="btn-action {type}" {title} on:click={handleClick}>
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
            {@html icons[type]}
        </svg>
    </a>
{:else}
    <button class="btn-action {type}" {title} on:click={handleClick}>
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
            {@html icons[type]}
        </svg>
    </button>
{/if}

<style>
    .btn-action {
        background: var(--color-bg-secondary, #f3f4f6);
        border: none;
        padding: 0.5rem;
        border-radius: 8px;
        cursor: pointer;
        font-size: 1rem;
        transition:
            background 0.2s,
            color 0.2s;
        text-decoration: none;
        display: inline-flex;
        align-items: center;
        justify-content: center;
        color: var(--color-text-secondary, #6b7280);
        width: 36px;
        height: 36px;
    }

    .btn-action:hover {
        background: var(--color-border, #e5e7eb);
        color: var(--color-text-primary, #111827);
    }

    .btn-action.delete:hover {
        background: var(--color-danger-bg, #fee2e2);
        color: var(--color-danger, #dc2626);
    }

    .btn-action svg {
        width: 18px;
        height: 18px;
    }
</style>
