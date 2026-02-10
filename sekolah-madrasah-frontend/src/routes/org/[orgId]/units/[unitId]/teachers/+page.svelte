<script lang="ts">
    import { page } from "$app/stores";
    import {
        api,
        Icon,
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
    let teachers: any[] = [];
    let isLoading = true;
    let error = "";
    let searchQuery = "";

    // Modals
    let showCreateModal = false;
    let showDeleteModal = false;
    let showResetPasswordModal = false;
    let selectedTeacher: any = null;
    let newPassword = "";

    // Pagination
    let currentPage = 1;
    let totalPages = 1;
    let totalItems = 0;
    let pageSize = 10;

    // Create form
    let newTeacher = {
        full_name: "",
        email: "",
        phone: "",
        nip: "",
        nuptk: "",
        education_level: "",
        education_major: "",
        employment_status: "PNS",
        join_date: "",
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

    async function loadTeachers() {
        try {
            isLoading = true;
            error = "";
            const response: any = await api.get(
                `/units/${unitId}/teachers?page=${currentPage}&limit=${pageSize}${searchQuery ? `&search=${searchQuery}` : ""}`,
            );
            if (response.data?.data) {
                teachers = response.data.data || [];
                totalItems = response.data.total || 0;
                totalPages = Math.ceil(totalItems / pageSize) || 1;
            } else {
                teachers = response.data || [];
                totalItems = teachers.length;
                totalPages = 1;
            }
        } catch (err: any) {
            error = err.message || "Gagal memuat data guru";
            teachers = [];
        } finally {
            isLoading = false;
        }
    }

    async function createTeacher() {
        try {
            const payload = {
                full_name: newTeacher.full_name,
                email: newTeacher.email || undefined,
                phone: newTeacher.phone || undefined,
                nip: newTeacher.nip || undefined,
                nuptk: newTeacher.nuptk || undefined,
                education_level: newTeacher.education_level || undefined,
                education_major: newTeacher.education_major || undefined,
                employment_status: newTeacher.employment_status,
                join_date: newTeacher.join_date || undefined,
            };

            const response: any = await api.post(
                `/units/${unitId}/teachers/with-user`,
                payload,
            );

            if (response.data?.password) {
                alert(
                    `Guru berhasil ditambahkan!\n\nEmail: ${response.data.email}\nPassword: ${response.data.password}\n\nSimpan password ini, tidak akan ditampilkan lagi.`,
                );
            }

            showCreateModal = false;
            resetForm();
            loadTeachers();
        } catch (err: any) {
            error = err.message || "Gagal menambah guru";
        }
    }

    async function deleteTeacher() {
        if (!selectedTeacher) return;
        try {
            await api.delete(`/units/${unitId}/teachers/${selectedTeacher.id}`);
            showDeleteModal = false;
            selectedTeacher = null;
            loadTeachers();
        } catch (err: any) {
            error = err.message || "Gagal menghapus guru";
        }
    }

    async function resetPassword() {
        if (!selectedTeacher) return;
        try {
            const chars =
                "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz23456789";
            let password = "";
            for (let i = 0; i < 10; i++) {
                password += chars.charAt(
                    Math.floor(Math.random() * chars.length),
                );
            }

            await api.put(`/users/${selectedTeacher.user_id}/reset-password`, {
                new_password: password,
            });

            newPassword = password;
        } catch (err: any) {
            error = err.message || "Gagal reset password";
        }
    }

    function resetForm() {
        newTeacher = {
            full_name: "",
            email: "",
            phone: "",
            nip: "",
            nuptk: "",
            education_level: "",
            education_major: "",
            employment_status: "PNS",
            join_date: "",
        };
    }

    function openDeleteModal(teacher: any) {
        selectedTeacher = teacher;
        showDeleteModal = true;
    }

    function openResetPasswordModal(teacher: any) {
        selectedTeacher = teacher;
        newPassword = "";
        showResetPasswordModal = true;
    }

    function handleSearch() {
        currentPage = 1;
        loadTeachers();
    }

    function handlePageChange(e: CustomEvent) {
        currentPage = e.detail.page;
        loadTeachers();
    }

    function getStatusBadgeClass(status: string) {
        switch (status?.toLowerCase()) {
            case "pns":
                return "badge-success";
            case "honorer":
                return "badge-warning";
            case "kontrak":
                return "badge-info";
            default:
                return "badge-default";
        }
    }

    onMount(async () => {
        await checkUserRole();
        loadTeachers();
    });
</script>

<svelte:head>
    <title>Daftar Guru - Sekolah Madrasah</title>
</svelte:head>

<div class="page-container">
    <!-- Page Header -->
    <PageHeader
        icon="users"
        title="Data Guru"
        subtitle="Kelola data guru dan tenaga pendidik"
        showAction={canManage}
        actionLabel="Tambah Guru"
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
        placeholder="Cari nama, NIP, atau NUPTK..."
        bind:value={searchQuery}
        {totalItems}
        itemLabel="guru"
        on:search={handleSearch}
    />

    <!-- Teachers Table -->
    <div class="table-container glass-card">
        {#if isLoading}
            <div class="loading-state">
                <div class="loader"></div>
                <p>Memuat data guru...</p>
            </div>
        {:else if teachers.length === 0}
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
                    >
                        <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
                        <circle cx="9" cy="7" r="4" />
                        <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
                        <path d="M16 3.13a4 4 0 0 1 0 7.75" />
                    </svg>
                </div>
                <h3>Belum Ada Data Guru</h3>
                <p>Mulai dengan menambahkan guru pertama</p>
                {#if canManage}
                    <button
                        class="btn-primary"
                        on:click={() => (showCreateModal = true)}
                    >
                        + Tambah Guru
                    </button>
                {/if}
            </div>
        {:else}
            <table class="data-table">
                <thead>
                    <tr>
                        <th>Foto</th>
                        <th>NIP / NUPTK</th>
                        <th>Nama Lengkap</th>
                        <th>Mata Pelajaran</th>
                        <th>Status</th>
                        <th>Pendidikan</th>
                        <th>Aksi</th>
                    </tr>
                </thead>
                <tbody>
                    {#each teachers as teacher}
                        <tr>
                            <td>
                                <div class="avatar">
                                    {teacher.user?.full_name?.charAt(0) || "G"}
                                </div>
                            </td>
                            <td>
                                <div class="nip-col">
                                    <span class="nip">{teacher.nip || "-"}</span
                                    >
                                    <span class="nuptk"
                                        >{teacher.nuptk || "-"}</span
                                    >
                                </div>
                            </td>
                            <td>
                                <a
                                    href="/org/{orgId}/units/{unitId}/teachers/{teacher.id}"
                                    class="teacher-name"
                                >
                                    {teacher.user?.full_name || "-"}
                                </a>
                            </td>
                            <td>
                                <div class="subjects-col">
                                    {#if teacher.teacher_subjects && teacher.teacher_subjects.length > 0}
                                        {#each teacher.teacher_subjects.slice(0, 2) as ts}
                                            <span class="subject-badge"
                                                >{ts.subject?.name || "-"}</span
                                            >
                                        {/each}
                                        {#if teacher.teacher_subjects.length > 2}
                                            <span class="subject-more"
                                                >+{teacher.teacher_subjects
                                                    .length - 2}</span
                                            >
                                        {/if}
                                    {:else}
                                        <span class="text-muted">-</span>
                                    {/if}
                                </div>
                            </td>
                            <td>
                                <span
                                    class="badge {getStatusBadgeClass(
                                        teacher.employment_status,
                                    )}"
                                >
                                    {teacher.employment_status || "-"}
                                </span>
                            </td>
                            <td>
                                <div class="education-col">
                                    <span>{teacher.education_level || "-"}</span
                                    >
                                    <span class="major"
                                        >{teacher.education_major || ""}</span
                                    >
                                </div>
                            </td>
                            <td>
                                <div class="actions">
                                    <ActionButton
                                        type="view"
                                        href="/org/{orgId}/units/{unitId}/teachers/{teacher.id}"
                                        title="Lihat Detail"
                                    />
                                    {#if canManage}
                                        <ActionButton
                                            type="reset"
                                            title="Reset Password"
                                            on:click={() =>
                                                openResetPasswordModal(teacher)}
                                        />
                                        <ActionButton
                                            type="delete"
                                            title="Hapus"
                                            on:click={() =>
                                                openDeleteModal(teacher)}
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
    title="Tambah Guru Baru"
    size="large"
    on:close={() => (showCreateModal = false)}
>
    <form on:submit|preventDefault={createTeacher}>
        <div class="form-section">
            <h3>Informasi Akun</h3>
            <div class="form-grid">
                <div class="form-group full-width">
                    <label for="full_name">Nama Lengkap *</label>
                    <input
                        type="text"
                        id="full_name"
                        bind:value={newTeacher.full_name}
                        required
                        placeholder="Nama lengkap guru"
                    />
                </div>
                <div class="form-group">
                    <label for="email">Email</label>
                    <input
                        type="email"
                        id="email"
                        bind:value={newTeacher.email}
                        placeholder="Email (opsional)"
                    />
                </div>
                <div class="form-group">
                    <label for="phone">No. Telepon</label>
                    <input
                        type="text"
                        id="phone"
                        bind:value={newTeacher.phone}
                        placeholder="08xxx"
                    />
                </div>
            </div>
        </div>

        <div class="form-section">
            <h3>Data Kepegawaian</h3>
            <div class="form-grid">
                <div class="form-group">
                    <label for="nip">NIP</label>
                    <input
                        type="text"
                        id="nip"
                        bind:value={newTeacher.nip}
                        placeholder="Nomor Induk Pegawai"
                    />
                </div>
                <div class="form-group">
                    <label for="nuptk">NUPTK</label>
                    <input
                        type="text"
                        id="nuptk"
                        bind:value={newTeacher.nuptk}
                        placeholder="Nomor Unik Pendidik"
                    />
                </div>
                <div class="form-group">
                    <label for="employment_status">Status Kepegawaian *</label>
                    <select
                        id="employment_status"
                        bind:value={newTeacher.employment_status}
                        required
                    >
                        <option value="PNS">PNS</option>
                        <option value="Honorer">Honorer</option>
                        <option value="Kontrak">Kontrak</option>
                        <option value="GTY">GTY (Guru Tetap Yayasan)</option>
                    </select>
                </div>
                <div class="form-group">
                    <label for="join_date">Tanggal Bergabung</label>
                    <input
                        type="date"
                        id="join_date"
                        bind:value={newTeacher.join_date}
                    />
                </div>
            </div>
        </div>

        <div class="form-section">
            <h3>Pendidikan</h3>
            <div class="form-grid">
                <div class="form-group">
                    <label for="education_level">Jenjang Pendidikan</label>
                    <select
                        id="education_level"
                        bind:value={newTeacher.education_level}
                    >
                        <option value="">Pilih Jenjang</option>
                        <option value="SMA">SMA/Sederajat</option>
                        <option value="D3">D3</option>
                        <option value="S1">S1</option>
                        <option value="S2">S2</option>
                        <option value="S3">S3</option>
                    </select>
                </div>
                <div class="form-group">
                    <label for="education_major">Jurusan/Prodi</label>
                    <input
                        type="text"
                        id="education_major"
                        bind:value={newTeacher.education_major}
                        placeholder="Contoh: Pendidikan Matematika"
                    />
                </div>
            </div>
        </div>

        <div class="form-info">
            <p>
                💡 Password akan otomatis dibuat berdasarkan NIP atau default
                "guru123456"
            </p>
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
        <button class="btn-primary" on:click={createTeacher}> Simpan </button>
    </svelte:fragment>
</Modal>

<!-- Delete Modal -->
<Modal
    show={showDeleteModal}
    title="Hapus Guru"
    size="small"
    on:close={() => (showDeleteModal = false)}
>
    <div class="delete-confirm">
        <p>Apakah Anda yakin ingin menghapus guru:</p>
        <p class="delete-name">{selectedTeacher?.user?.full_name || "-"}</p>
        <p class="delete-warning">Tindakan ini tidak dapat dibatalkan!</p>
    </div>
    <svelte:fragment slot="actions">
        <button
            class="btn-secondary"
            on:click={() => (showDeleteModal = false)}
        >
            Batal
        </button>
        <button class="btn-danger" on:click={deleteTeacher}> Hapus </button>
    </svelte:fragment>
</Modal>

<!-- Reset Password Modal -->
<Modal
    show={showResetPasswordModal}
    title="Reset Password"
    size="small"
    on:close={() => (showResetPasswordModal = false)}
>
    <div class="reset-password-content">
        {#if newPassword}
            <p>
                Password baru untuk <strong
                    >{selectedTeacher?.user?.full_name}</strong
                >:
            </p>
            <div class="password-display">
                <code>{newPassword}</code>
                <button
                    class="btn-copy"
                    on:click={() => navigator.clipboard.writeText(newPassword)}
                >
                    📋 Salin
                </button>
            </div>
            <p class="delete-warning">
                Simpan password ini, tidak akan ditampilkan lagi!
            </p>
        {:else}
            <p>Reset password untuk:</p>
            <p class="delete-name">{selectedTeacher?.user?.full_name || "-"}</p>
            <p class="text-muted">
                Password baru akan dihasilkan secara otomatis
            </p>
        {/if}
    </div>
    <svelte:fragment slot="actions">
        <button
            class="btn-secondary"
            on:click={() => (showResetPasswordModal = false)}
        >
            {newPassword ? "Tutup" : "Batal"}
        </button>
        {#if !newPassword}
            <button class="btn-primary" on:click={resetPassword}>
                Reset Password
            </button>
        {/if}
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

    /* Avatar */
    .avatar {
        width: 40px;
        height: 40px;
        border-radius: 50%;
        background: #00ced1;
        color: white;
        display: flex;
        align-items: center;
        justify-content: center;
        font-weight: 600;
    }

    /* Teacher-specific styles */
    .nip-col {
        display: flex;
        flex-direction: column;
        gap: 0.25rem;
    }

    .nip {
        font-family: monospace;
        font-size: 0.875rem;
    }

    .nuptk {
        font-family: monospace;
        font-size: 0.75rem;
        color: #9ca3af;
    }

    .teacher-name {
        font-weight: 500;
        color: #00ced1;
        text-decoration: none;
    }

    .teacher-name:hover {
        text-decoration: underline;
    }

    .subjects-col {
        display: flex;
        flex-wrap: wrap;
        gap: 0.25rem;
    }

    .subject-badge {
        background: #e0f7f7;
        color: #00ced1;
        padding: 0.125rem 0.5rem;
        border-radius: 4px;
        font-size: 0.75rem;
    }

    .subject-more {
        background: #f3f4f6;
        color: #6b7280;
        padding: 0.125rem 0.5rem;
        border-radius: 4px;
        font-size: 0.75rem;
    }

    .education-col {
        display: flex;
        flex-direction: column;
        gap: 0.125rem;
    }

    .major {
        font-size: 0.75rem;
        color: #9ca3af;
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

    .badge-success {
        background: #d1fae5;
        color: #10b981;
    }

    .badge-warning {
        background: #fef3c7;
        color: #f59e0b;
    }

    .badge-info {
        background: #dbeafe;
        color: #3b82f6;
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
    .form-section {
        margin-bottom: 1.5rem;
    }

    .form-section h3 {
        font-size: 1rem;
        font-weight: 600;
        color: #374151;
        margin: 0 0 1rem;
        padding-bottom: 0.5rem;
        border-bottom: 1px solid #e5e7eb;
    }

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
    .form-group select {
        padding: 0.75rem;
        border: 1px solid #e5e7eb;
        border-radius: 8px;
        font-size: 1rem;
        background: #ffffff;
        transition: border-color 0.2s;
    }

    .form-group input:focus,
    .form-group select:focus {
        outline: none;
        border-color: #00ced1;
    }

    .form-info {
        background: #e0f7f7;
        padding: 0.75rem 1rem;
        border-radius: 8px;
        margin-bottom: 1rem;
    }

    .form-info p {
        margin: 0;
        font-size: 0.875rem;
        color: #00ced1;
    }

    /* Delete/Reset Confirm */
    .delete-confirm,
    .reset-password-content {
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

    .password-display {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
        background: #f3f4f6;
        padding: 1rem;
        border-radius: 8px;
        margin: 1rem 0;
    }

    .password-display code {
        font-size: 1.25rem;
        font-weight: 600;
        color: #111827;
    }

    .btn-copy {
        background: #00ced1;
        color: white;
        border: none;
        padding: 0.5rem 1rem;
        border-radius: 6px;
        cursor: pointer;
        font-size: 0.875rem;
    }

    .btn-copy:hover {
        background: #00b5b8;
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
