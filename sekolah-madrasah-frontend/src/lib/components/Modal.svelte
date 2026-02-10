<script lang="ts">
    /**
     * Modal Component
     * ===============
     * Reusable modal dialog with header, content, and actions slots.
     *
     * @example
     * <Modal show={showModal} title="Tambah Data" on:close={() => showModal = false}>
     *   <form>...</form>
     *   <svelte:fragment slot="actions">
     *     <button class="btn-secondary">Batal</button>
     *     <button class="btn-primary">Simpan</button>
     *   </svelte:fragment>
     * </Modal>
     */

    import { createEventDispatcher } from "svelte";
    import { fade, scale } from "svelte/transition";

    export let show: boolean = false;
    export let title: string = "Modal Title";
    export let size: "default" | "small" | "large" = "default";

    const dispatch = createEventDispatcher();

    function handleClose() {
        dispatch("close");
    }

    function handleOverlayClick(e: MouseEvent) {
        if (e.target === e.currentTarget) {
            handleClose();
        }
    }

    function handleKeydown(e: KeyboardEvent) {
        if (e.key === "Escape") {
            handleClose();
        }
    }
</script>

<svelte:window on:keydown={handleKeydown} />

{#if show}
    <div
        class="modal-overlay"
        on:click={handleOverlayClick}
        on:keypress={() => {}}
        transition:fade={{ duration: 150 }}
    >
        <div
            class="modal-content glass-card modal-{size}"
            on:click|stopPropagation
            on:keypress={() => {}}
            transition:scale={{ duration: 150, start: 0.95 }}
        >
            <div class="modal-header">
                <h2>{title}</h2>
                <button class="btn-close" on:click={handleClose}>×</button>
            </div>

            <div class="modal-body">
                <slot />
            </div>

            {#if $$slots.actions}
                <div class="modal-actions">
                    <slot name="actions" />
                </div>
            {/if}
        </div>
    </div>
{/if}

<style>
    .modal-overlay {
        position: fixed;
        inset: 0;
        background: rgba(0, 0, 0, 0.5);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 1000;
        padding: 1rem;
    }

    .modal-content {
        width: 100%;
        max-height: 90vh;
        overflow-y: auto;
        padding: 0;
        background: var(--color-bg-card, #ffffff);
        border-radius: 16px;
        box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
    }

    .modal-default {
        max-width: 600px;
    }

    .modal-small {
        max-width: 400px;
    }

    .modal-large {
        max-width: 800px;
    }

    .modal-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1.5rem;
        border-bottom: 1px solid var(--color-border, #e5e7eb);
    }

    .modal-header h2 {
        font-size: 1.25rem;
        font-weight: 600;
        color: var(--color-text-primary, #111827);
        margin: 0;
    }

    .btn-close {
        width: 32px;
        height: 32px;
        border: none;
        background: var(--color-bg-secondary, #f3f4f6);
        border-radius: 8px;
        font-size: 1.25rem;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        color: var(--color-text-secondary, #6b7280);
        transition: background 0.2s;
    }

    .btn-close:hover {
        background: var(--color-border, #e5e7eb);
    }

    .modal-body {
        padding: 1.5rem;
    }

    .modal-actions {
        display: flex;
        justify-content: flex-end;
        gap: 0.75rem;
        padding: 1.5rem;
        border-top: 1px solid var(--color-border, #e5e7eb);
    }

    @media (max-width: 768px) {
        .modal-content {
            max-width: 100%;
            max-height: 100%;
            border-radius: 0;
        }
    }
</style>
