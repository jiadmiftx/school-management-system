<!--
  DataTable Component
  ===================
  Reusable data table with loading and empty states.
  
  Usage:
  <DataTable
    columns={[
      { key: 'name', label: 'Nama' },
      { key: 'email', label: 'Email' },
    ]}
    data={items}
    isLoading={isLoading}
    emptyTitle="Belum Ada Data"
    emptyMessage="Mulai dengan menambahkan data baru"
  >
    <svelte:fragment slot="row" let:item let:index>
      <td>{item.name}</td>
      <td>{item.email}</td>
      <td>
        <div class="actions">
          <ActionButton type="edit" on:click={() => edit(item)} />
          <ActionButton type="delete" on:click={() => remove(item)} />
        </div>
      </td>
    </svelte:fragment>
    
    <svelte:fragment slot="empty-action">
      <button class="btn-primary" on:click={create}>+ Tambah</button>
    </svelte:fragment>
  </DataTable>
-->
<script lang="ts">
    interface Column {
        key: string;
        label: string;
        width?: string;
    }

    export let columns: Column[] = [];
    export let data: any[] = [];
    export let isLoading: boolean = false;
    export let emptyIcon: string = "file";
    export let emptyTitle: string = "Belum Ada Data";
    export let emptyMessage: string = "Data belum tersedia";
</script>

<div class="table-container glass-card">
    {#if isLoading}
        <div class="loading-state">
            <div class="loader"></div>
            <p>Memuat data...</p>
        </div>
    {:else if data.length === 0}
        <div class="empty-state">
            <div class="empty-icon">
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="64"
                    height="64"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="#9ca3af"
                    stroke-width="1.5"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                >
                    <path
                        d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z"
                    />
                    <polyline points="14,2 14,8 20,8" />
                    <line x1="12" y1="18" x2="12" y2="12" />
                    <line x1="9" y1="15" x2="15" y2="15" />
                </svg>
            </div>
            <h3>{emptyTitle}</h3>
            <p>{emptyMessage}</p>
            <slot name="empty-action" />
        </div>
    {:else}
        <table class="data-table">
            <thead>
                <tr>
                    {#each columns as column}
                        <th
                            style={column.width ? `width: ${column.width}` : ""}
                        >
                            {column.label}
                        </th>
                    {/each}
                </tr>
            </thead>
            <tbody>
                {#each data as item, index}
                    <tr>
                        <slot name="row" {item} {index} />
                    </tr>
                {/each}
            </tbody>
        </table>
    {/if}
</div>

<style>
    .table-container {
        overflow: hidden;
    }

    .loading-state,
    .empty-state {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 4rem 2rem;
        color: var(--color-text-secondary, #6b7280);
    }

    .loader {
        width: 48px;
        height: 48px;
        border: 4px solid var(--color-border, #e5e7eb);
        border-top-color: var(--color-primary, #00ced1);
        border-radius: 50%;
        animation: spin 1s linear infinite;
    }

    @keyframes spin {
        to {
            transform: rotate(360deg);
        }
    }

    .empty-icon {
        margin-bottom: 1rem;
    }

    .empty-state h3 {
        font-size: 1.25rem;
        color: var(--color-text-primary, #374151);
        margin: 0 0 0.5rem;
    }

    .empty-state p {
        margin: 0 0 1.5rem;
    }

    .data-table {
        width: 100%;
        border-collapse: collapse;
    }

    .data-table th,
    .data-table :global(td) {
        padding: 1rem 1.5rem;
        text-align: left;
        border-bottom: 1px solid var(--color-border-light, #f3f4f6);
    }

    .data-table th {
        background: var(--color-bg-hover, #f9fafb);
        font-weight: 600;
        color: var(--color-text-primary, #374151);
        font-size: 0.875rem;
        text-transform: uppercase;
        letter-spacing: 0.025em;
    }

    .data-table tbody tr:hover {
        background: var(--color-bg-hover, #f9fafb);
    }
</style>
