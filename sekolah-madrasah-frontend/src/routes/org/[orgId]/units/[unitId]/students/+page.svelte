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
    let students: any[] = [];
    let isLoading = true;
    let error = "";
    let searchQuery = "";

    // Modals
    let showCreateModal = false;
    let showDeleteModal = false;
    let showResetPasswordModal = false;
    let selectedStudent: any = null;
    let newPassword = "";

    // Pagination
    let currentPage = 1;
    let totalPages = 1;
    let totalItems = 0;
    let pageSize = 10;

    // Create form
    let newStudent = {
        nis: "",
        nisn: "",
        full_name: "",
        email: "",
        birth_place: "",
        birth_date: "",
        gender: "",
        religion: "",
        address: "",
        father_name: "",
        mother_name: "",
        guardian_name: "",
        parent_phone: "",
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

    async function loadStudents() {
        try {
            isLoading = true;
            error = "";
            const response: any = await api.get(
                `/units/${unitId}/students?page=${currentPage}&limit=${pageSize}${searchQuery ? `&search=${searchQuery}` : ""}`,
            );
            students = response.data || [];
            if (response.paginate) {
                totalItems = response.paginate.total_data || 0;
                totalPages = response.paginate.total_pages || 1;
            } else {
                totalItems = students.length;
                totalPages = 1;
            }
        } catch (err: any) {
            error = err.message || "Gagal memuat data siswa";
            students = [];
        } finally {
            isLoading = false;
        }
    }

    async function createStudent() {
        try {
            const payload = {
                full_name: newStudent.full_name,
                email: newStudent.email || undefined,
                nis: newStudent.nis,
                nisn: newStudent.nisn || undefined,
                birth_place: newStudent.birth_place || undefined,
                birth_date: newStudent.birth_date,
                gender: newStudent.gender || undefined,
                religion: newStudent.religion || undefined,
                address: newStudent.address || undefined,
                father_name: newStudent.father_name || undefined,
                mother_name: newStudent.mother_name || undefined,
                guardian_name: newStudent.guardian_name || undefined,
                parent_phone: newStudent.parent_phone || undefined,
            };

            const response: any = await api.post(
                `/units/${unitId}/students/with-user`,
                payload,
            );

            if (response.data?.password) {
                alert(
                    `Siswa berhasil ditambahkan!\n\nPassword: ${response.data.password}\n\nSimpan password ini, tidak akan ditampilkan lagi.`,
                );
            }

            showCreateModal = false;
            resetForm();
            loadStudents();
        } catch (err: any) {
            error = err.message || "Gagal menambah siswa";
        }
    }

    async function deleteStudent() {
        if (!selectedStudent) return;
        try {
            await api.delete(`/units/${unitId}/students/${selectedStudent.id}`);
            showDeleteModal = false;
            selectedStudent = null;
            loadStudents();
        } catch (err: any) {
            error = err.message || "Gagal menghapus siswa";
        }
    }

    async function resetPassword() {
        if (!selectedStudent) return;
        try {
            const chars =
                "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz23456789";
            let password = "";
            for (let i = 0; i < 10; i++) {
                password += chars.charAt(
                    Math.floor(Math.random() * chars.length),
                );
            }

            await api.put(`/users/${selectedStudent.user_id}/reset-password`, {
                new_password: password,
            });

            newPassword = password;
        } catch (err: any) {
            error = err.message || "Gagal reset password";
        }
    }

    function resetForm() {
        newStudent = {
            nis: "",
            nisn: "",
            full_name: "",
            email: "",
            birth_place: "",
            birth_date: "",
            gender: "",
            religion: "",
            address: "",
            father_name: "",
            mother_name: "",
            guardian_name: "",
            parent_phone: "",
        };
    }

    function openDeleteModal(student: any) {
        selectedStudent = student;
        showDeleteModal = true;
    }

    function openResetPasswordModal(student: any) {
        selectedStudent = student;
        newPassword = "";
        showResetPasswordModal = true;
    }

    function handleSearch() {
        currentPage = 1;
        loadStudents();
    }

    function handlePageChange(e: CustomEvent) {
        currentPage = e.detail.page;
        loadStudents();
    }

    function formatDate(dateString: string) {
        if (!dateString) return "-";
        return new Date(dateString).toLocaleDateString("id-ID", {
            day: "numeric",
            month: "long",
            year: "numeric",
        });
    }

    onMount(async () => {
        await checkUserRole();
        loadStudents();
    });
</script>

<svelte:head>
    <title>Daftar Siswa - Sekolah Madrasah</title>
</svelte:head>

<div class="page-container">
    <!-- Page Header -->
    <PageHeader
        icon="graduation"
        title="Data Siswa"
        subtitle="Kelola data siswa dan profil lengkapnya"
        showAction={canManage}
        actionLabel="Tambah Siswa"
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
        placeholder="Cari nama, NIS, atau NISN..."
        bind:value={searchQuery}
        {totalItems}
        itemLabel="siswa"
        on:search={handleSearch}
    />

    <!-- Students Table -->
    <div class="table-container glass-card">
        {#if isLoading}
            <div class="loading-state">
                <div class="loader"></div>
                <p>Memuat data siswa...</p>
            </div>
        {:else if students.length === 0}
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
                        <path d="M22 10v6M2 10l10-5 10 5-10 5z" />
                        <path d="M6 12v5c3 3 9 3 12 0v-5" />
                    </svg>
                </div>
                <h3>Belum Ada Data Siswa</h3>
                <p>Mulai dengan menambahkan siswa pertama</p>
                {#if canManage}
                    <button
                        class="btn-primary"
                        on:click={() => (showCreateModal = true)}
                    >
                        + Tambah Siswa
                    </button>
                {/if}
            </div>
        {:else}
            <table class="data-table">
                <thead>
                    <tr>
                        <th>Foto</th>
                        <th>NIS / NISN</th>
                        <th>Nama Lengkap</th>
                        <th>Jenis Kelamin</th>
                        <th>Tanggal Lahir</th>
                        <th>Aksi</th>
                    </tr>
                </thead>
                <tbody>
                    {#each students as student}
                        <tr>
                            <td>
                                <div class="avatar">
                                    {student.user?.full_name?.charAt(0) || "S"}
                                </div>
                            </td>
                            <td>
                                <div class="nis-col">
                                    <span class="nis">{student.nis || "-"}</span
                                    >
                                    <span class="nisn"
                                        >{student.nisn || "-"}</span
                                    >
                                </div>
                            </td>
                            <td>
                                <a
                                    href="/org/{orgId}/units/{unitId}/students/{student.id}"
                                    class="student-name"
                                >
                                    {student.user?.full_name || "-"}
                                </a>
                            </td>
                            <td>
                                <span
                                    class="gender-badge"
                                    class:male={student.gender === "L"}
                                    class:female={student.gender === "P"}
                                >
                                    {student.gender === "L"
                                        ? "Laki-laki"
                                        : student.gender === "P"
                                          ? "Perempuan"
                                          : "-"}
                                </span>
                            </td>
                            <td>{formatDate(student.birth_date)}</td>
                            <td>
                                <div class="actions">
                                    <ActionButton
                                        type="view"
                                        href="/org/{orgId}/units/{unitId}/students/{student.id}"
                                        title="Lihat Detail"
                                    />
                                    {#if canManage}
                                        <ActionButton
                                            type="reset"
                                            title="Reset Password"
                                            on:click={() =>
                                                openResetPasswordModal(student)}
                                        />
                                        <ActionButton
                                            type="delete"
                                            title="Hapus"
                                            on:click={() =>
                                                openDeleteModal(student)}
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
    title="Tambah Siswa Baru"
    size="large"
    on:close={() => (showCreateModal = false)}
>
    <form on:submit|preventDefault={createStudent}>
        <div class="form-section">
            <h3>Data Identitas</h3>
            <div class="form-grid">
                <div class="form-group">
                    <label for="nis">NIS <span class="required">*</span></label>
                    <input
                        type="text"
                        id="nis"
                        bind:value={newStudent.nis}
                        placeholder="Nomor Induk Siswa"
                        required
                    />
                </div>
                <div class="form-group">
                    <label for="nisn">NISN</label>
                    <input
                        type="text"
                        id="nisn"
                        bind:value={newStudent.nisn}
                        placeholder="Nomor Induk Siswa Nasional"
                    />
                </div>
                <div class="form-group full-width">
                    <label for="full_name"
                        >Nama Lengkap <span class="required">*</span></label
                    >
                    <input
                        type="text"
                        id="full_name"
                        bind:value={newStudent.full_name}
                        placeholder="Nama lengkap siswa"
                        required
                    />
                </div>
                <div class="form-group full-width">
                    <label for="email">Email (opsional)</label>
                    <input
                        type="email"
                        id="email"
                        bind:value={newStudent.email}
                        placeholder="email@contoh.com (jika kosong akan auto-generate)"
                    />
                </div>
            </div>
        </div>

        <div class="form-section">
            <h3>Data Pribadi</h3>
            <div class="form-grid">
                <div class="form-group">
                    <label for="birth_place">Tempat Lahir</label>
                    <input
                        type="text"
                        id="birth_place"
                        bind:value={newStudent.birth_place}
                        placeholder="Kota kelahiran"
                    />
                </div>
                <div class="form-group">
                    <label for="birth_date"
                        >Tanggal Lahir <span class="required">*</span></label
                    >
                    <input
                        type="date"
                        id="birth_date"
                        bind:value={newStudent.birth_date}
                        required
                    />
                </div>
                <div class="form-group">
                    <label for="gender">Jenis Kelamin</label>
                    <select id="gender" bind:value={newStudent.gender}>
                        <option value="">Pilih...</option>
                        <option value="L">Laki-laki</option>
                        <option value="P">Perempuan</option>
                    </select>
                </div>
                <div class="form-group">
                    <label for="religion">Agama</label>
                    <select id="religion" bind:value={newStudent.religion}>
                        <option value="">Pilih...</option>
                        <option value="Islam">Islam</option>
                        <option value="Kristen">Kristen</option>
                        <option value="Katolik">Katolik</option>
                        <option value="Hindu">Hindu</option>
                        <option value="Buddha">Buddha</option>
                        <option value="Konghucu">Konghucu</option>
                    </select>
                </div>
                <div class="form-group full-width">
                    <label for="address">Alamat</label>
                    <textarea
                        id="address"
                        bind:value={newStudent.address}
                        placeholder="Alamat lengkap"
                        rows="2"
                    ></textarea>
                </div>
            </div>
        </div>

        <div class="form-section">
            <h3>Data Orang Tua/Wali</h3>
            <div class="form-grid">
                <div class="form-group">
                    <label for="father_name">Nama Ayah</label>
                    <input
                        type="text"
                        id="father_name"
                        bind:value={newStudent.father_name}
                        placeholder="Nama lengkap ayah"
                    />
                </div>
                <div class="form-group">
                    <label for="mother_name">Nama Ibu</label>
                    <input
                        type="text"
                        id="mother_name"
                        bind:value={newStudent.mother_name}
                        placeholder="Nama lengkap ibu"
                    />
                </div>
                <div class="form-group">
                    <label for="guardian_name">Nama Wali</label>
                    <input
                        type="text"
                        id="guardian_name"
                        bind:value={newStudent.guardian_name}
                        placeholder="Nama wali (jika ada)"
                    />
                </div>
                <div class="form-group">
                    <label for="parent_phone">No. Telepon Orang Tua</label>
                    <input
                        type="text"
                        id="parent_phone"
                        bind:value={newStudent.parent_phone}
                        placeholder="08xxx"
                    />
                </div>
            </div>
        </div>

        <div class="form-info">
            <p>
                💡 Password akan otomatis dibuat berdasarkan NIS atau default
                "siswa123456"
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
        <button class="btn-primary" on:click={createStudent}> Simpan </button>
    </svelte:fragment>
</Modal>

<!-- Delete Modal -->
<Modal
    show={showDeleteModal}
    title="Hapus Siswa"
    size="small"
    on:close={() => (showDeleteModal = false)}
>
    <div class="delete-confirm">
        <p>Apakah Anda yakin ingin menghapus siswa:</p>
        <p class="delete-name">{selectedStudent?.user?.full_name || "-"}</p>
        <p class="delete-warning">Tindakan ini tidak dapat dibatalkan!</p>
    </div>
    <svelte:fragment slot="actions">
        <button
            class="btn-secondary"
            on:click={() => (showDeleteModal = false)}
        >
            Batal
        </button>
        <button class="btn-danger" on:click={deleteStudent}> Hapus </button>
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
                    >{selectedStudent?.user?.full_name}</strong
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
            <p class="delete-name">{selectedStudent?.user?.full_name || "-"}</p>
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

    /* Student-specific styles */
    .nis-col {
        display: flex;
        flex-direction: column;
        gap: 0.25rem;
    }

    .nis {
        font-family: monospace;
        font-size: 0.875rem;
    }

    .nisn {
        font-family: monospace;
        font-size: 0.75rem;
        color: #9ca3af;
    }

    .student-name {
        font-weight: 500;
        color: #00ced1;
        text-decoration: none;
    }

    .student-name:hover {
        text-decoration: underline;
    }

    .gender-badge {
        display: inline-block;
        padding: 0.25rem 0.75rem;
        border-radius: 9999px;
        font-size: 0.75rem;
        font-weight: 500;
        background: #f3f4f6;
        color: #6b7280;
    }

    .gender-badge.male {
        background: #dbeafe;
        color: #3b82f6;
    }

    .gender-badge.female {
        background: #fce7f3;
        color: #ec4899;
    }

    .text-muted {
        color: #9ca3af;
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

    .required {
        color: #dc2626;
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
