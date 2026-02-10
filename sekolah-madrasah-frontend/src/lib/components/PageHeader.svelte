<!--
  PageHeader Component
  ====================
  Reusable page header with icon, title, subtitle, and action button.
  
  Usage:
  <PageHeader
    icon="users"
    title="Data Guru"
    subtitle="Kelola data guru dan tenaga pendidik"
    showAction={canManage}
    actionLabel="Tambah Guru"
    on:action={() => showCreateModal = true}
  />
-->
<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import { Icon } from "$lib";

    export let icon: string = "file";
    export let title: string = "Page Title";
    export let subtitle: string = "";
    export let showAction: boolean = true;
    export let actionLabel: string = "Tambah";
    export let actionIcon: string = "+";

    const dispatch = createEventDispatcher();

    function handleAction() {
        dispatch("action");
    }
</script>

<div class="page-header glass-card">
    <div class="header-content">
        <div class="header-icon">
            <Icon name={icon} size={28} color="white" />
        </div>
        <div class="header-text">
            <h1>{title}</h1>
            {#if subtitle}
                <p class="subtitle">{subtitle}</p>
            {/if}
        </div>
    </div>

    {#if showAction}
        <button class="btn-primary" on:click={handleAction}>
            <span class="icon">{actionIcon}</span>
            {actionLabel}
        </button>
    {/if}
</div>

<style>
    .page-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1.5rem 2rem;
        margin-bottom: 1.5rem;
        background: #ffffff;
        border: 1px solid #e5e7eb;
        border-radius: 16px;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
    }

    .header-content {
        display: flex;
        align-items: center;
        gap: 1rem;
    }

    .header-icon {
        width: 56px;
        height: 56px;
        background: var(--color-primary, #00ced1);
        padding: 0.75rem;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .header-text h1 {
        color: var(--color-text-primary, #111827);
        font-size: 1.75rem;
        font-weight: 700;
        margin: 0;
    }

    .subtitle {
        color: var(--color-text-secondary, #6b7280);
        margin: 0.25rem 0 0;
        font-size: 0.95rem;
    }

    .btn-primary {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        background: var(--color-primary, #00ced1);
        color: white;
        border: none;
        padding: 0.75rem 1.5rem;
        border-radius: 12px;
        font-size: 1rem;
        font-weight: 600;
        cursor: pointer;
        transition:
            transform 0.2s,
            box-shadow 0.2s,
            background 0.2s;
    }

    .btn-primary:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 20px rgba(0, 206, 209, 0.4);
        background: var(--color-primary-hover, #00b5b8);
    }

    .btn-primary .icon {
        font-size: 1.25rem;
    }
</style>
