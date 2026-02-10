<script lang="ts">
    import { page } from "$app/stores";
    import { goto } from "$app/navigation";
    import { api } from "$lib";
    import { onMount } from "svelte";
    import { isSuperAdmin } from "$core/stores/auth";

    // Get IDs from URL
    $: orgId = $page.params.orgId;
    $: unitId = $page.params.unitId;
    $: teacherId = $page.params.teacherId;

    // State
    let teacher: any = null;
    let isLoading = true;
    let error = "";
    let isEditing = false;
    let isSaving = false;

    // Subjects
    let allSubjects: any[] = [];
    let showAssignSubjectModal = false;
    let selectedSubjectToAssign = "";
    let isPrimarySubject = false;

    // Role check
    let canManage = false;
    let userRole = "";

    // Edit form
    let editForm = {
        nip: "",
        nuptk: "",
        education_level: "",
        education_major: "",
        employment_status: "",
        join_date: "",
    };

    async function checkUserRole() {
        try {
            // Superadmin has full access
            if ($isSuperAdmin) {
                canManage = true;
                return;
            }

            const membershipRes: any = await api.get("/users/me/memberships");
            const memberships = membershipRes.data?.unit_memberships || [];
            const membership = memberships.find(
                (m: any) => m.unit_id === unitId,
            );
            if (membership) {
                userRole = membership.role;
                canManage =
                    userRole === "admin" ||
                    userRole === "staff" ||
                    userRole === "owner";
            }
        } catch (err) {
            console.error("Failed to check role:", err);
        }
    }

    async function loadTeacher() {
        try {
            isLoading = true;
            const response: any = await api.get(
                `/units/${unitId}/teachers/${teacherId}`,
            );
            if (response.data) {
                teacher = response.data;
                populateEditForm();
            }
        } catch (err: any) {
            error = err.message || "Gagal memuat data guru";
        } finally {
            isLoading = false;
        }
    }

    function populateEditForm() {
        if (!teacher) return;
        editForm = {
            nip: teacher.nip || "",
            nuptk: teacher.nuptk || "",
            education_level: teacher.education_level || "",
            education_major: teacher.education_major || "",
            employment_status: teacher.employment_status || "",
            join_date: teacher.join_date ? teacher.join_date.split("T")[0] : "",
        };
    }

    async function saveChanges() {
        try {
            isSaving = true;
            await api.put(`/units/${unitId}/teachers/${teacherId}`, editForm);
            await loadTeacher();
            isEditing = false;
        } catch (err: any) {
            error = err.message || "Gagal menyimpan perubahan";
        } finally {
            isSaving = false;
        }
    }

    function cancelEdit() {
        isEditing = false;
        populateEditForm();
    }

    async function loadSubjects() {
        try {
            const response: any = await api.get(
                `/units/${unitId}/subjects?limit=100`,
            );
            if (response.data?.data) {
                allSubjects = response.data.data || [];
            } else {
                allSubjects = response.data || [];
            }
        } catch (err) {
            console.error("Failed to load subjects:", err);
        }
    }

    async function assignSubject() {
        if (!selectedSubjectToAssign) return;
        try {
            await api.post(`/subjects/${selectedSubjectToAssign}/teachers`, {
                teacher_profile_id: teacherId,
                is_primary: isPrimarySubject,
            });
            showAssignSubjectModal = false;
            selectedSubjectToAssign = "";
            isPrimarySubject = false;
            await loadTeacher();
        } catch (err: any) {
            error = err.message || "Gagal menambahkan mata pelajaran";
        }
    }

    async function removeSubject(subjectId: string) {
        if (!confirm("Apakah Anda yakin ingin menghapus mapel ini dari guru?"))
            return;
        try {
            await api.delete(`/subjects/${subjectId}/teachers/${teacherId}`);
            await loadTeacher();
        } catch (err: any) {
            error = err.message || "Gagal menghapus mata pelajaran";
        }
    }

    function formatDate(dateString: string) {
        if (!dateString) return "-";
        return new Date(dateString).toLocaleDateString("id-ID", {
            day: "numeric",
            month: "long",
            year: "numeric",
        });
    }

    function getStatusBadgeClass(status: string) {
        switch (status?.toLowerCase()) {
            case "pns":
                return "status-pns";
            case "honorer":
                return "status-honorer";
            case "kontrak":
                return "status-kontrak";
            default:
                return "";
        }
    }

    function goBack() {
        goto(`/org/${orgId}/units/${unitId}/teachers`);
    }

    onMount(async () => {
        await checkUserRole();
        loadTeacher();
        loadSubjects();
    });
</script>

<svelte:head>
    <title>{teacher?.user?.full_name || "Detail Guru"} - Sekolah Madrasah</title
    >
</svelte:head>

<div class="teacher-detail-container">
    <!-- Back Button -->
    <button class="btn-back" on:click={goBack}>
        ← Kembali ke Daftar Guru
    </button>

    {#if error}
        <div class="error-message glass-card">
            <span>⚠️ {error}</span>
            <button on:click={() => (error = "")}>×</button>
        </div>
    {/if}

    {#if isLoading}
        <div class="loading-state glass-card">
            <div class="loader"></div>
            <p>Memuat data guru...</p>
        </div>
    {:else if teacher}
        <!-- Profile Header -->
        <div class="profile-header glass-card">
            <div class="profile-avatar">
                <span>{teacher.user?.full_name?.charAt(0) || "G"}</span>
            </div>
            <div class="profile-info">
                <h1>{teacher.user?.full_name || "-"}</h1>
                <div class="profile-badges">
                    {#if teacher.nip}
                        <span class="badge nip">NIP: {teacher.nip}</span>
                    {/if}
                    {#if teacher.nuptk}
                        <span class="badge nuptk">NUPTK: {teacher.nuptk}</span>
                    {/if}
                    {#if teacher.employment_status}
                        <span
                            class="badge status {getStatusBadgeClass(
                                teacher.employment_status,
                            )}"
                        >
                            {teacher.employment_status}
                        </span>
                    {/if}
                </div>
            </div>
            {#if canManage}
                <div class="profile-actions">
                    {#if isEditing}
                        <button class="btn-cancel" on:click={cancelEdit}
                            >Batal</button
                        >
                        <button
                            class="btn-save"
                            on:click={saveChanges}
                            disabled={isSaving}
                        >
                            {isSaving ? "Menyimpan..." : "Simpan"}
                        </button>
                    {:else}
                        <button
                            class="btn-edit"
                            on:click={() => (isEditing = true)}
                        >
                            Edit Profil
                        </button>
                    {/if}
                </div>
            {/if}
        </div>

        <!-- Content -->
        <div class="content-section glass-card">
            <h2>Data Kepegawaian</h2>
            <div class="info-grid">
                <div class="info-item">
                    <label>NIP</label>
                    {#if isEditing}
                        <input
                            type="text"
                            bind:value={editForm.nip}
                            placeholder="Nomor Induk Pegawai"
                        />
                    {:else}
                        <span>{teacher.nip || "-"}</span>
                    {/if}
                </div>
                <div class="info-item">
                    <label>NUPTK</label>
                    {#if isEditing}
                        <input
                            type="text"
                            bind:value={editForm.nuptk}
                            placeholder="Nomor Unik Pendidik"
                        />
                    {:else}
                        <span>{teacher.nuptk || "-"}</span>
                    {/if}
                </div>
                <div class="info-item">
                    <label>Status Kepegawaian</label>
                    {#if isEditing}
                        <select bind:value={editForm.employment_status}>
                            <option value="">Pilih Status</option>
                            <option value="PNS">PNS</option>
                            <option value="Honorer">Honorer</option>
                            <option value="Kontrak">Kontrak</option>
                            <option value="GTY">GTY (Guru Tetap Yayasan)</option
                            >
                        </select>
                    {:else}
                        <span
                            class="badge status {getStatusBadgeClass(
                                teacher.employment_status,
                            )}"
                        >
                            {teacher.employment_status || "-"}
                        </span>
                    {/if}
                </div>
                <div class="info-item">
                    <label>Tanggal Bergabung</label>
                    {#if isEditing}
                        <input type="date" bind:value={editForm.join_date} />
                    {:else}
                        <span>{formatDate(teacher.join_date)}</span>
                    {/if}
                </div>
            </div>
        </div>

        <div class="content-section glass-card">
            <h2>Pendidikan</h2>
            <div class="info-grid">
                <div class="info-item">
                    <label>Jenjang Pendidikan</label>
                    {#if isEditing}
                        <select bind:value={editForm.education_level}>
                            <option value="">Pilih Jenjang</option>
                            <option value="SMA">SMA/Sederajat</option>
                            <option value="D3">D3</option>
                            <option value="S1">S1</option>
                            <option value="S2">S2</option>
                            <option value="S3">S3</option>
                        </select>
                    {:else}
                        <span>{teacher.education_level || "-"}</span>
                    {/if}
                </div>
                <div class="info-item">
                    <label>Jurusan/Prodi</label>
                    {#if isEditing}
                        <input
                            type="text"
                            bind:value={editForm.education_major}
                            placeholder="Contoh: Pendidikan Matematika"
                        />
                    {:else}
                        <span>{teacher.education_major || "-"}</span>
                    {/if}
                </div>
            </div>
        </div>

        <div class="content-section glass-card">
            <div class="section-header-row">
                <h2>Mata Pelajaran</h2>
                {#if canManage}
                    <button
                        class="btn-add-subject"
                        on:click={() => (showAssignSubjectModal = true)}
                    >
                        + Tambah Mapel
                    </button>
                {/if}
            </div>
            <div class="subjects-display">
                {#if teacher.teacher_subjects && teacher.teacher_subjects.length > 0}
                    {#each teacher.teacher_subjects as ts}
                        <div class="subject-card">
                            <span class="subject-name"
                                >{ts.subject?.name || "-"}</span
                            >
                            {#if ts.is_primary}
                                <span class="subject-primary">Utama</span>
                            {/if}
                            {#if canManage}
                                <button
                                    class="btn-remove-subject"
                                    title="Hapus mapel"
                                    on:click={() =>
                                        removeSubject(ts.subject?.id)}
                                >
                                    ×
                                </button>
                            {/if}
                        </div>
                    {/each}
                {:else}
                    <p class="no-subjects">
                        Belum ada mata pelajaran yang ditugaskan
                    </p>
                {/if}
            </div>
        </div>

        <div class="content-section glass-card">
            <h2>Informasi Akun</h2>
            <div class="info-grid">
                <div class="info-item">
                    <label>Email</label>
                    <span>{teacher.user?.email || "-"}</span>
                </div>
                <div class="info-item">
                    <label>No. Telepon</label>
                    <span>{teacher.user?.phone || "-"}</span>
                </div>
            </div>
        </div>
    {:else}
        <div class="empty-state glass-card">
            <h3>Data Guru Tidak Ditemukan</h3>
            <p>Guru yang Anda cari mungkin sudah dihapus atau tidak ada.</p>
            <button class="btn-back-main" on:click={goBack}>
                Kembali ke Daftar Guru
            </button>
        </div>
    {/if}
</div>

<!-- Assign Subject Modal -->
{#if showAssignSubjectModal}
    <div
        class="modal-overlay"
        on:click={() => (showAssignSubjectModal = false)}
        on:keypress={() => {}}
    >
        <div
            class="modal-content glass-card"
            on:click|stopPropagation
            on:keypress={() => {}}
        >
            <div class="modal-header">
                <h2>Tambah Mata Pelajaran</h2>
                <button
                    class="btn-close"
                    on:click={() => (showAssignSubjectModal = false)}>×</button
                >
            </div>
            <form on:submit|preventDefault={assignSubject}>
                <div class="form-group">
                    <label for="subject">Pilih Mata Pelajaran *</label>
                    <select
                        id="subject"
                        bind:value={selectedSubjectToAssign}
                        required
                    >
                        <option value="">-- Pilih Mapel --</option>
                        {#each allSubjects.filter((s) => !teacher?.teacher_subjects?.some((ts) => ts.subject?.id === s.id)) as subject}
                            <option value={subject.id}
                                >{subject.name} ({subject.code})</option
                            >
                        {/each}
                    </select>
                </div>
                <div class="form-group checkbox-group">
                    <label class="checkbox-label">
                        <input
                            type="checkbox"
                            bind:checked={isPrimarySubject}
                        />
                        <span>Jadikan mapel utama</span>
                    </label>
                </div>
                <div class="modal-actions">
                    <button
                        type="button"
                        class="btn-cancel"
                        on:click={() => (showAssignSubjectModal = false)}
                    >
                        Batal
                    </button>
                    <button type="submit" class="btn-submit"> Simpan </button>
                </div>
            </form>
        </div>
    </div>
{/if}

<style>
    /* Container */
    .teacher-detail-container {
        padding: 2rem;
        min-height: 100vh;
        background: #f5f5f7;
    }

    /* Card */
    .glass-card {
        background: #ffffff;
        border: 1px solid #e5e7eb;
        border-radius: 16px;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
    }

    /* Back Button */
    .btn-back {
        background: #ffffff;
        border: 1px solid #e5e7eb;
        color: #374151;
        padding: 0.75rem 1.25rem;
        border-radius: 10px;
        cursor: pointer;
        font-size: 0.95rem;
        margin-bottom: 1.5rem;
        transition: all 0.2s;
    }

    .btn-back:hover {
        background: #f3f4f6;
        border-color: #d1d5db;
    }

    /* Error */
    .error-message {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1rem 1.5rem;
        background: #fef2f2;
        border: 1px solid #fee2e2;
        margin-bottom: 1rem;
        color: #dc2626;
    }

    .error-message button {
        background: none;
        border: none;
        font-size: 1.5rem;
        cursor: pointer;
        color: #dc2626;
    }

    /* Loading */
    .loading-state {
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

    /* Profile Header */
    .profile-header {
        display: flex;
        align-items: center;
        gap: 1.5rem;
        padding: 2rem;
        margin-bottom: 1.5rem;
    }

    .profile-avatar {
        width: 80px;
        height: 80px;
        background: #00ced1;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 2rem;
        font-weight: 700;
        color: white;
        border: 3px solid #ffffff;
        box-shadow: 0 4px 12px rgba(0, 206, 209, 0.3);
        flex-shrink: 0;
    }

    .profile-avatar span {
        font-size: 2rem;
        font-weight: 700;
        color: white;
    }

    .profile-info {
        flex: 1;
    }

    .profile-info h1 {
        font-size: 1.5rem;
        font-weight: 700;
        color: #111827;
        margin: 0 0 0.5rem;
    }

    .profile-badges {
        display: flex;
        flex-wrap: wrap;
        gap: 0.5rem;
    }

    .badge {
        padding: 0.25rem 0.75rem;
        border-radius: 20px;
        font-size: 0.75rem;
        font-weight: 500;
    }

    .badge.nip {
        background: #dbeafe;
        color: #1e40af;
    }

    .badge.nuptk {
        background: #f3e8ff;
        color: #7c3aed;
    }

    .status-pns {
        background: #d1fae5;
        color: #065f46;
    }

    .status-honorer {
        background: #fef3c7;
        color: #92400e;
    }

    .status-kontrak {
        background: #e0e7ff;
        color: #4338ca;
    }

    .profile-actions {
        display: flex;
        gap: 0.5rem;
    }

    .btn-edit,
    .btn-save {
        padding: 0.75rem 1.5rem;
        background: #00ced1;
        color: white;
        border: none;
        border-radius: 10px;
        font-weight: 600;
        cursor: pointer;
        transition: all 0.2s;
    }

    .btn-edit:hover,
    .btn-save:hover {
        background: #00b5b8;
    }

    .btn-save:disabled {
        opacity: 0.7;
        cursor: not-allowed;
    }

    .btn-cancel {
        padding: 0.75rem 1.5rem;
        background: #f3f4f6;
        color: #374151;
        border: none;
        border-radius: 10px;
        font-weight: 500;
        cursor: pointer;
    }

    /* Content Sections */
    .content-section {
        padding: 1.5rem 2rem;
        margin-bottom: 1.5rem;
    }

    .content-section h2 {
        font-size: 1rem;
        font-weight: 600;
        color: #6b7280;
        text-transform: uppercase;
        letter-spacing: 0.05em;
        margin: 0 0 1.5rem;
        padding-bottom: 0.75rem;
        border-bottom: 1px solid #f3f4f6;
    }

    .info-grid {
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        gap: 1.5rem;
    }

    .info-item {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .info-item label {
        font-size: 0.875rem;
        font-weight: 500;
        color: #6b7280;
    }

    .info-item span {
        font-size: 1rem;
        color: #111827;
    }

    .info-item input,
    .info-item select {
        padding: 0.75rem 1rem;
        border: 1px solid #e5e7eb;
        border-radius: 10px;
        font-size: 0.95rem;
        transition: all 0.2s;
    }

    .info-item input:focus,
    .info-item select:focus {
        outline: none;
        border-color: #00ced1;
        box-shadow: 0 0 0 3px rgba(0, 206, 209, 0.1);
    }

    /* Subjects Display */
    .subjects-display {
        display: flex;
        flex-wrap: wrap;
        gap: 0.75rem;
    }

    .subject-card {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.75rem 1rem;
        background: #e0f2fe;
        border-radius: 10px;
        border: 1px solid #bae6fd;
    }

    .subject-name {
        font-weight: 500;
        color: #0369a1;
        font-size: 0.95rem;
    }

    .subject-primary {
        padding: 0.2rem 0.5rem;
        background: #0369a1;
        color: white;
        border-radius: 6px;
        font-size: 0.7rem;
        font-weight: 600;
        text-transform: uppercase;
    }

    .no-subjects {
        color: #9ca3af;
        font-style: italic;
        margin: 0;
    }

    /* Section Header Row */
    .section-header-row {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1rem;
    }

    .section-header-row h2 {
        margin: 0;
    }

    .btn-add-subject {
        padding: 0.5rem 1rem;
        background: #00ced1;
        color: white;
        border: none;
        border-radius: 8px;
        font-size: 0.85rem;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.2s;
    }

    .btn-add-subject:hover {
        background: #00b5b8;
    }

    .btn-remove-subject {
        width: 24px;
        height: 24px;
        background: #fee2e2;
        color: #dc2626;
        border: none;
        border-radius: 6px;
        font-size: 1rem;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: all 0.2s;
    }

    .btn-remove-subject:hover {
        background: #fecaca;
    }

    /* Modal */
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
        max-width: 400px;
        padding: 0;
    }

    .modal-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1.25rem 1.5rem;
        border-bottom: 1px solid #e5e7eb;
    }

    .modal-header h2 {
        font-size: 1.1rem;
        font-weight: 600;
        color: #111827;
        margin: 0;
    }

    .btn-close {
        width: 32px;
        height: 32px;
        border: none;
        background: #f3f4f6;
        border-radius: 8px;
        font-size: 1.25rem;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .form-group {
        padding: 0 1.5rem;
        margin-top: 1rem;
    }

    .form-group label {
        display: block;
        font-size: 0.875rem;
        font-weight: 500;
        color: #374151;
        margin-bottom: 0.5rem;
    }

    .form-group select {
        width: 100%;
        padding: 0.75rem 1rem;
        border: 1px solid #e5e7eb;
        border-radius: 10px;
        font-size: 0.95rem;
    }

    .form-group select:focus {
        outline: none;
        border-color: #00ced1;
        box-shadow: 0 0 0 3px rgba(0, 206, 209, 0.1);
    }

    .checkbox-group {
        padding-bottom: 1rem;
    }

    .checkbox-label {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        cursor: pointer;
    }

    .checkbox-label input {
        width: 18px;
        height: 18px;
        cursor: pointer;
    }

    .modal-actions {
        display: flex;
        justify-content: flex-end;
        gap: 0.75rem;
        padding: 1rem 1.5rem;
        border-top: 1px solid #e5e7eb;
        background: #f9fafb;
        border-radius: 0 0 16px 16px;
    }

    .btn-cancel {
        padding: 0.75rem 1.5rem;
        background: #f3f4f6;
        color: #374151;
        border: none;
        border-radius: 10px;
        font-weight: 500;
        cursor: pointer;
    }

    .btn-submit {
        padding: 0.75rem 1.5rem;
        background: #00ced1;
        color: white;
        border: none;
        border-radius: 10px;
        font-weight: 600;
        cursor: pointer;
    }

    .btn-submit:hover {
        background: #00b5b8;
    }

    /* Empty State */
    .empty-state {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 4rem 2rem;
        text-align: center;
    }

    .empty-state h3 {
        font-size: 1.25rem;
        color: #374151;
        margin: 0 0 0.5rem;
    }

    .empty-state p {
        color: #6b7280;
        margin: 0 0 1.5rem;
    }

    .btn-back-main {
        padding: 0.75rem 1.5rem;
        background: #00ced1;
        color: white;
        border: none;
        border-radius: 10px;
        font-weight: 600;
        cursor: pointer;
    }

    /* Responsive */
    @media (max-width: 768px) {
        .teacher-detail-container {
            padding: 0 1rem 1rem;
        }

        .profile-header {
            flex-direction: column;
            text-align: center;
            padding: 1.5rem;
        }

        .profile-badges {
            justify-content: center;
        }

        .profile-actions {
            width: 100%;
            justify-content: center;
        }

        .info-grid {
            grid-template-columns: 1fr;
        }
    }
</style>
