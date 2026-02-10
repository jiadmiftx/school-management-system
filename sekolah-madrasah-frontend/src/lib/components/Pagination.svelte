<!--
  Pagination Component
  ====================
  Reusable pagination controls.
  
  Usage:
  <Pagination
    currentPage={currentPage}
    totalPages={totalPages}
    on:change={(e) => { currentPage = e.detail.page; loadData(); }}
  />
-->
<script lang="ts">
    import { createEventDispatcher } from "svelte";

    export let currentPage: number = 1;
    export let totalPages: number = 1;

    const dispatch = createEventDispatcher();

    function goToPrevious() {
        if (currentPage > 1) {
            dispatch("change", { page: currentPage - 1 });
        }
    }

    function goToNext() {
        if (currentPage < totalPages) {
            dispatch("change", { page: currentPage + 1 });
        }
    }
</script>

{#if totalPages > 1}
    <div class="pagination">
        <button
            class="btn-page"
            disabled={currentPage === 1}
            on:click={goToPrevious}
        >
            ← Prev
        </button>
        <span class="page-info">
            Halaman {currentPage} dari {totalPages}
        </span>
        <button
            class="btn-page"
            disabled={currentPage === totalPages}
            on:click={goToNext}
        >
            Next →
        </button>
    </div>
{/if}

<style>
    .pagination {
        display: flex;
        justify-content: center;
        align-items: center;
        gap: 1rem;
        padding: 1.5rem 0;
        border-top: 1px solid var(--color-border-light, #f3f4f6);
    }

    .btn-page {
        padding: 0.5rem 1rem;
        background: var(--color-bg-secondary, #f3f4f6);
        border: none;
        border-radius: 8px;
        cursor: pointer;
        font-size: 0.875rem;
        color: var(--color-text-primary, #111827);
        transition: background 0.2s;
    }

    .btn-page:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }

    .btn-page:not(:disabled):hover {
        background: var(--color-border, #e5e7eb);
    }

    .page-info {
        color: var(--color-text-secondary, #6b7280);
        font-size: 0.875rem;
    }
</style>
