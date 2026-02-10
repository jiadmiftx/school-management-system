<script lang="ts">
    import { page } from "$app/stores";
    import {
        api,
        PageHeader,
        SearchBar,
        ActionButton,
        Modal,
        Pagination,
    } from "$lib";
    import { onMount } from "svelte";

    // Get IDs from URL
    $: orgId = $page.params.orgId;
    $: unitId = $page.params.unitId;

    // State
    let subjects: any[] = [];
    let isLoading = true;
    let error = "";
    let searchQuery = "";

    // Modals
    let showCreateModal = false;
    let showEditModal = false;
    let showDeleteModal = false;
    let selectedSubject: any = null;

    // Pagination
    let currentPage = 1;
    let totalPages = 1;
    let totalItems = 0;
    let pageSize = 10;

    // Form
    let formData = {
        name: "",
        code: "",
        category: "",
        description: "",
    };

    // Role check
    let canManage = true;
    let userRole = "";

    async function checkUserRole() {
        try {
            const membershipRes: any = await api.get("/users/me/memberships");
            const memberships = membershipRes.data?.unit_memberships || [];
            const membership = memberships.find(
                (m: any) => m.unit_id === unitId,
            );
            if (membership) {
                userRole = membership.role;
                canManage =
                    userRole === "owner" ||
                    userRole === "admin" ||
                    userRole === "staff";
            }
        } catch (err) {
            console.error("Failed to check role:", err);
        }
    }

    async function loadSubjects() {
        try {
            isLoading = true;
            error = "";
            const response: any = await api.get(
                `/units/${unitId}/subjects?page=${currentPage}&limit=${pageSize}${searchQuery ? `&search=${searchQuery}` : ""}`,
            );
            if (response.data?.data) {
                subjects = response.data.data || [];
                totalItems = response.data.total || 0;
                totalPages = Math.ceil(totalItems / pageSize) || 1;
            } else {
                subjects = response.data || [];
                totalItems = subjects.length;
                totalPages = 1;
            }
        } catch (err: any) {
            error = err.message || "Gagal memuat data mata pelajaran";
            subjects = [];
        } finally {
            isLoading = false;
        }
    }

    async function createSubject() {
        try {
            const payload = {
                name: formData.name,
                code: formData.code,
                category: formData.category || undefined,
                description: formData.description || undefined,
            };

            await api.post(`/units/${unitId}/subjects`, payload);

            showCreateModal = false;
            resetForm();
            loadSubjects();
        } catch (err: any) {
            error = err.message || "Gagal menambah mata pelajaran";
        }
    }

    async function updateSubject() {
        if (!selectedSubject) return;
        try {
            const payload = {
                name: formData.name,
                code: formData.code,
                category: formData.category || undefined,
                description: formData.description || undefined,
            };

            await api.put(`/subjects/${selectedSubject.id}`, payload);

            showEditModal = false;
            selectedSubject = null;
            resetForm();
            loadSubjects();
        } catch (err: any) {
            error = err.message || "Gagal memperbarui mata pelajaran";
        }
    }

    async function deleteSubject() {
        if (!selectedSubject) return;
        try {
            await api.delete(`/subjects/${selectedSubject.id}`);
            showDeleteModal = false;
            selectedSubject = null;
            loadSubjects();
        } catch (err: any) {
            error = err.message || "Gagal menghapus mata pelajaran";
        }
    }

    function resetForm() {
        formData = {
            name: "",
            code: "",
            category: "",
            description: "",
        };
    }

    function openEditModal(subject: any) {
        selectedSubject = subject;
        formData = {
            name: subject.name || "",
            code: subject.code || "",
            category: subject.category || "",
            description: subject.description || "",
        };
        showEditModal = true;
    }

    function openDeleteModal(subject: any) {
        selectedSubject = subject;
        showDeleteModal = true;
    }

    function handleSearch() {
        currentPage = 1;
        loadSubjects();
    }

    function handlePageChange(e: CustomEvent) {
        currentPage = e.detail.page;
        loadSubjects();
    }

    function getCategoryBadgeClass(category: string) {
        switch (category?.toLowerCase()) {
            case "umum":
                return "badge-info";
            case "jurusan":
                return "badge-success";
            case "mulok":
                return "badge-warning";
            default:
                return "badge-default";
        }
    }

    onMount(async () => {
        await checkUserRole();
        loadSubjects();
    });
</script>

<svelte:head>
    <title>Mata Pelajaran - Sekolah Madrasah</title>
</svelte:head>

<div class="page-container">
    <!-- Page Header -->
    <PageHeader
        icon="book"
        title="Mata Pelajaran"
        subtitle="Kelola data mata pelajaran"
        showAction={canManage}
        actionLabel="Tambah Mapel"
        on:action={() => (showCreateModal = true)}
    />

    <!-- Error Message -->
    {#if error}
        <div class="error-message glass-card">
            <span>⚠️ {error}</span>
            <button on:click={() => (error = "")}>×</button>
        </div>
    {/if}

    <!-- Search Bar -->
    <SearchBar
        placeholder="Cari nama atau kode mapel..."
        bind:value={searchQuery}
        {totalItems}
        itemLabel="mapel"
        on:search={handleSearch}
    />

    <!-- Subjects Table -->
    <div class="table-container glass-card">
        {#if isLoading}
            <div class="loading-state">
                <div class="loader"></div>
                <p>Memuat data mata pelajaran...</p>
            </div>
        {:else if subjects.length === 0}
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
                        <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20" />
                        <path
                            d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"
                        />
                    </svg>
                </div>
                <h3>Belum Ada Data Mata Pelajaran</h3>
                <p>Mulai dengan menambahkan mata pelajaran pertama</p>
                {#if canManage}
                    <button
                        class="btn-primary"
                        on:click={() => (showCreateModal = true)}
                    >
                        + Tambah Mapel
                    </button>
                {/if}
            </div>
        {:else}
            <table class="data-table">
                <thead>
                    <tr>
                        <th>Kode</th>
                        <th>Nama Mata Pelajaran</th>
                        <th>Kategori</th>
                        <th>Deskripsi</th>
                        <th>Aksi</th>
                    </tr>
                </thead>
                <tbody>
                    {#each subjects as subject}
                        <tr>
                            <td>
                                <span class="subject-code"
                                    >{subject.code || "-"}</span
                                >
                            </td>
                            <td>
                                <span class="subject-name"
                                    >{subject.name || "-"}</span
                                >
                            </td>
                            <td>
                                {#if subject.category}
                                    <span
                                        class="badge {getCategoryBadgeClass(
                                            subject.category,
                                        )}"
                                    >
                                        {subject.category}
                                    </span>
                                {:else}
                                    <span class="text-muted">-</span>
                                {/if}
                            </td>
                            <td>
                                <span class="description"
                                    >{subject.description || "-"}</span
                                >
                            </td>
                            <td>
                                <div class="actions">
                                    {#if canManage}
                                        <ActionButton
                                            type="edit"
                                            title="Edit"
                                            on:click={() =>
                                                openEditModal(subject)}
                                        />
                                        <ActionButton
                                            type="delete"
                                            title="Hapus"
                                            on:click={() =>
                                                openDeleteModal(subject)}
                                        />
                                    {/if}
                                </div>
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>

            <!-- Pagination -->
            <Pagination
                {currentPage}
                {totalPages}
                on:change={handlePageChange}
            />
        {/if}
    </div>
</div>

<!-- Create Modal -->
<Modal
    show={showCreateModal}
    title="Tambah Mata Pelajaran"
    on:close={() => (showCreateModal = false)}
>
    <form on:submit|preventDefault={createSubject}>
        <div class="form-grid">
            <div class="form-group">
                <label for="name">Nama Mata Pelajaran *</label>
                <input
                    type="text"
                    id="name"
                    bind:value={formData.name}
                    placeholder="Contoh: Matematika"
                    required
                />
            </div>
            <div class="form-group">
                <label for="code">Kode Mapel *</label>
                <input
                    type="text"
                    id="code"
                    bind:value={formData.code}
                    placeholder="Contoh: MTK"
                    required
                />
            </div>
            <div class="form-group">
                <label for="category">Kategori</label>
                <select id="category" bind:value={formData.category}>
                    <option value="">Pilih Kategori</option>
                    <option value="Umum">Umum</option>
                    <option value="Jurusan">Jurusan</option>
                    <option value="Mulok">Muatan Lokal</option>
                </select>
            </div>
            <div class="form-group full-width">
                <label for="description">Deskripsi</label>
                <textarea
                    id="description"
                    bind:value={formData.description}
                    placeholder="Deskripsi mata pelajaran (opsional)"
                    rows="3"
                ></textarea>
            </div>
        </div>
    </form>
    <svelte:fragment slot="actions">
        <button
            type="button"
            class="btn-secondary"
            on:click={() => (showCreateModal = false)}
        >
            Batal
        </button>
        <button type="submit" class="btn-primary" on:click={createSubject}>
            Simpan
        </button>
    </svelte:fragment>
</Modal>

<!-- Edit Modal -->
<Modal
    show={showEditModal}
    title="Edit Mata Pelajaran"
    on:close={() => (showEditModal = false)}
>
    <form on:submit|preventDefault={updateSubject}>
        <div class="form-grid">
            <div class="form-group">
                <label for="edit-name">Nama Mata Pelajaran *</label>
                <input
                    type="text"
                    id="edit-name"
                    bind:value={formData.name}
                    placeholder="Contoh: Matematika"
                    required
                />
            </div>
            <div class="form-group">
                <label for="edit-code">Kode Mapel *</label>
                <input
                    type="text"
                    id="edit-code"
                    bind:value={formData.code}
                    placeholder="Contoh: MTK"
                    required
                />
            </div>
            <div class="form-group">
                <label for="edit-category">Kategori</label>
                <select id="edit-category" bind:value={formData.category}>
                    <option value="">Pilih Kategori</option>
                    <option value="Umum">Umum</option>
                    <option value="Jurusan">Jurusan</option>
                    <option value="Mulok">Muatan Lokal</option>
                </select>
            </div>
            <div class="form-group full-width">
                <label for="edit-description">Deskripsi</label>
                <textarea
                    id="edit-description"
                    bind:value={formData.description}
                    placeholder="Deskripsi mata pelajaran (opsional)"
                    rows="3"
                ></textarea>
            </div>
        </div>
    </form>
    <svelte:fragment slot="actions">
        <button
            type="button"
            class="btn-secondary"
            on:click={() => (showEditModal = false)}
        >
            Batal
        </button>
        <button type="submit" class="btn-primary" on:click={updateSubject}>
            Simpan Perubahan
        </button>
    </svelte:fragment>
</Modal>

<!-- Delete Modal -->
<Modal
    show={showDeleteModal}
    title="Hapus Mata Pelajaran"
    size="small"
    on:close={() => (showDeleteModal = false)}
>
    <div class="delete-confirm">
        <p>Apakah Anda yakin ingin menghapus mata pelajaran:</p>
        <p class="delete-name">
            {selectedSubject?.name} ({selectedSubject?.code})
        </p>
        <p class="delete-warning">Tindakan ini tidak dapat dibatalkan.</p>
    </div>
    <svelte:fragment slot="actions">
        <button
            class="btn-secondary"
            on:click={() => (showDeleteModal = false)}
        >
            Batal
        </button>
        <button class="btn-danger" on:click={deleteSubject}> Hapus </button>
    </svelte:fragment>
</Modal>

<style>
    /* Page Container */
    .page-container {
        padding: 2rem;
        min-height: 100vh;
        background: #f5f5f7;
    }

    /* Glass Card */
    .glass-card {
        background: #ffffff;
        border: 1px solid #e5e7eb;
        border-radius: 16px;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
    }

    /* Error Message */
    .error-message {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1rem 1.5rem;
        margin-bottom: 1.5rem;
        background: #fee2e2;
        border: 1px solid #fecaca;
        color: #dc2626;
    }

    .error-message button {
        background: none;
        border: none;
        font-size: 1.5rem;
        cursor: pointer;
        color: #dc2626;
    }

    /* Table Container */
    .table-container {
        overflow: hidden;
    }

    /* Loading & Empty States */
    .loading-state,
    .empty-state {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 4rem 2rem;
        color: #6b7280;
    }

    .loader {
        width: 48px;
        height: 48px;
        border: 4px solid #e5e7eb;
        border-top-color: #00ced1;
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
        color: #374151;
        margin: 0 0 0.5rem;
    }

    .empty-state p {
        margin: 0 0 1.5rem;
    }

    /* Data Table */
    .data-table {
        width: 100%;
        border-collapse: collapse;
    }

    .data-table th,
    .data-table td {
        padding: 1rem 1.5rem;
        text-align: left;
        border-bottom: 1px solid #f3f4f6;
    }

    .data-table th {
        background: #f9fafb;
        font-weight: 600;
        color: #374151;
        font-size: 0.875rem;
        text-transform: uppercase;
        letter-spacing: 0.025em;
    }

    .data-table tbody tr:hover {
        background: #f9fafb;
    }

    .actions {
        display: flex;
        gap: 0.5rem;
    }

    /* Subject-specific styles */
    .subject-code {
        font-family: monospace;
        background: #f3f4f6;
        padding: 0.25rem 0.5rem;
        border-radius: 4px;
        font-size: 0.875rem;
    }

    .subject-name {
        font-weight: 500;
        color: #111827;
    }

    .description {
        color: #6b7280;
        font-size: 0.875rem;
        max-width: 200px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    }

    .text-muted {
        color: #9ca3af;
    }

    /* Badges */
    .badge {
        display: inline-block;
        padding: 0.25rem 0.75rem;
        border-radius: 9999px;
        font-size: 0.75rem;
        font-weight: 600;
        text-transform: uppercase;
    }

    .badge-info {
        background: #dbeafe;
        color: #3b82f6;
    }

    .badge-success {
        background: #d1fae5;
        color: #10b981;
    }

    .badge-warning {
        background: #fef3c7;
        color: #f59e0b;
    }

    .badge-default {
        background: #f3f4f6;
        color: #6b7280;
    }

    /* Buttons */
    .btn-primary {
        display: inline-flex;
        align-items: center;
        gap: 0.5rem;
        background: #00ced1;
        color: white;
        border: none;
        padding: 0.75rem 1.5rem;
        border-radius: 12px;
        font-size: 1rem;
        font-weight: 600;
        cursor: pointer;
        transition:
            transform 0.2s,
            box-shadow 0.2s;
    }

    .btn-primary:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 20px rgba(0, 206, 209, 0.4);
    }

    .btn-secondary {
        background: #f3f4f6;
        color: #374151;
        border: none;
        padding: 0.75rem 1.5rem;
        border-radius: 12px;
        font-weight: 500;
        cursor: pointer;
        transition: background 0.2s;
    }

    .btn-secondary:hover {
        background: #e5e7eb;
    }

    .btn-danger {
        background: #dc2626;
        color: white;
        border: none;
        padding: 0.75rem 1.5rem;
        border-radius: 12px;
        font-weight: 600;
        cursor: pointer;
        transition: background 0.2s;
    }

    .btn-danger:hover {
        background: #b91c1c;
    }

    /* Form */
    .form-grid {
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        gap: 1rem;
    }

    .form-group {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .form-group.full-width {
        grid-column: span 2;
    }

    .form-group label {
        font-weight: 500;
        color: #374151;
        font-size: 0.875rem;
    }

    .form-group input,
    .form-group select,
    .form-group textarea {
        padding: 0.75rem;
        border: 1px solid #e5e7eb;
        border-radius: 8px;
        font-size: 1rem;
        background: #ffffff;
        transition: border-color 0.2s;
    }

    .form-group input:focus,
    .form-group select:focus,
    .form-group textarea:focus {
        outline: none;
        border-color: #00ced1;
    }

    /* Delete Confirm */
    .delete-confirm {
        text-align: center;
    }

    .delete-name {
        font-weight: 600;
        color: #111827;
        font-size: 1.125rem;
        margin: 0.5rem 0;
    }

    .delete-warning {
        color: #dc2626;
        font-size: 0.875rem;
    }

    @media (max-width: 768px) {
        .page-container {
            padding: 1rem;
        }

        .form-grid {
            grid-template-columns: 1fr;
        }

        .form-group.full-width {
            grid-column: span 1;
        }
    }
</style>
