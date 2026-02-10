<script lang="ts">
    import { page } from "$app/stores";
    import { goto } from "$app/navigation";
    import { api } from "$lib";
    import { onMount } from "svelte";

    // Get IDs from URL
    $: orgId = $page.params.orgId;
    $: unitId = $page.params.unitId;
    $: studentId = $page.params.studentId;

    // State
    let student: any = null;
    let isLoading = true;
    let error = "";
    let isEditing = false;
    let isSaving = false;

    // Role check
    let canManage = false;
    let userRole = "";

    // Edit form
    let editForm = {
        nis: "",
        nisn: "",
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

    // Active tab
    let activeTab = "biodata";

    async function checkUserRole() {
        try {
            const membershipRes: any = await api.get("/users/me/memberships");
            const memberships = membershipRes.data?.unit_memberships || [];
            const membership = memberships.find(
                (m: any) => m.unit_id === unitId,
            );
            if (membership) {
                userRole = membership.role;
                canManage = userRole === "admin" || userRole === "staff";
            }
        } catch (err) {
            console.error("Failed to check role:", err);
        }
    }

    async function loadStudent() {
        try {
            isLoading = true;
            const response = await api.get(
                `/units/${unitId}/students/${studentId}`,
            );
            if (response.data) {
                student = response.data;
                populateEditForm();
            }
        } catch (err: any) {
            error = err.message || "Gagal memuat data siswa";
        } finally {
            isLoading = false;
        }
    }

    function populateEditForm() {
        if (!student) return;
        editForm = {
            nis: student.nis || "",
            nisn: student.nisn || "",
            birth_place: student.birth_place || "",
            birth_date: student.birth_date
                ? student.birth_date.split("T")[0]
                : "",
            gender: student.gender || "",
            religion: student.religion || "",
            address: student.address || "",
            father_name: student.father_name || "",
            mother_name: student.mother_name || "",
            guardian_name: student.guardian_name || "",
            parent_phone: student.parent_phone || "",
        };
    }

    async function saveChanges() {
        try {
            isSaving = true;
            await api.put(`/units/${unitId}/students/${studentId}`, editForm);
            await loadStudent();
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

    function formatDate(dateString: string) {
        if (!dateString) return "-";
        return new Date(dateString).toLocaleDateString("id-ID", {
            day: "numeric",
            month: "long",
            year: "numeric",
        });
    }

    function goBack() {
        goto(`/org/${orgId}/units/${unitId}/students`);
    }

    onMount(async () => {
        await checkUserRole();
        loadStudent();
    });
</script>

<svelte:head>
    <title
        >{student?.user?.full_name || "Detail Siswa"} - Sekolah Madrasah</title
    >
</svelte:head>

<div class="student-detail-container">
    <!-- Back Button -->
    <button class="btn-back" on:click={goBack}>
        ← Kembali ke Daftar Siswa
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
            <p>Memuat data siswa...</p>
        </div>
    {:else if student}
        <!-- Profile Header -->
        <div class="profile-header glass-card">
            <div class="profile-avatar">
                <span>{student.user?.full_name?.charAt(0) || "S"}</span>
            </div>
            <div class="profile-info">
                <h1>{student.user?.full_name || "-"}</h1>
                <div class="profile-badges">
                    <span class="badge nis">NIS: {student.nis || "-"}</span>
                    <span class="badge nisn">NISN: {student.nisn || "-"}</span>
                    {#if student.gender}
                        <span
                            class="badge gender"
                            class:male={student.gender === "L"}
                            class:female={student.gender === "P"}
                        >
                            {student.gender === "L" ? "Laki-laki" : "Perempuan"}
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

        <!-- Tabs -->
        <div class="tabs glass-card">
            <button
                class="tab"
                class:active={activeTab === "biodata"}
                on:click={() => (activeTab = "biodata")}
            >
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
                    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
                    <circle cx="12" cy="7" r="4" />
                </svg>
                Data Pribadi
            </button>
            <button
                class="tab"
                class:active={activeTab === "parents"}
                on:click={() => (activeTab = "parents")}
            >
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
                    <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
                    <circle cx="9" cy="7" r="4" />
                    <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
                    <path d="M16 3.13a4 4 0 0 1 0 7.75" />
                </svg>
                Orang Tua
            </button>
            <button
                class="tab"
                class:active={activeTab === "enrollment"}
                on:click={() => (activeTab = "enrollment")}
            >
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
                    <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20" />
                    <path
                        d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"
                    />
                </svg>
                Riwayat Kelas
            </button>
        </div>

        <!-- Tab Content -->
        <div class="tab-content glass-card">
            {#if activeTab === "biodata"}
                <div class="info-grid">
                    <div class="info-item">
                        <label>NIS</label>
                        {#if isEditing}
                            <input type="text" bind:value={editForm.nis} />
                        {:else}
                            <span>{student.nis || "-"}</span>
                        {/if}
                    </div>
                    <div class="info-item">
                        <label>NISN</label>
                        {#if isEditing}
                            <input type="text" bind:value={editForm.nisn} />
                        {:else}
                            <span>{student.nisn || "-"}</span>
                        {/if}
                    </div>
                    <div class="info-item">
                        <label>Tempat Lahir</label>
                        {#if isEditing}
                            <input
                                type="text"
                                bind:value={editForm.birth_place}
                            />
                        {:else}
                            <span>{student.birth_place || "-"}</span>
                        {/if}
                    </div>
                    <div class="info-item">
                        <label>Tanggal Lahir</label>
                        {#if isEditing}
                            <input
                                type="date"
                                bind:value={editForm.birth_date}
                            />
                        {:else}
                            <span>{formatDate(student.birth_date)}</span>
                        {/if}
                    </div>
                    <div class="info-item">
                        <label>Jenis Kelamin</label>
                        {#if isEditing}
                            <select bind:value={editForm.gender}>
                                <option value="">Pilih...</option>
                                <option value="L">Laki-laki</option>
                                <option value="P">Perempuan</option>
                            </select>
                        {:else}
                            <span
                                >{student.gender === "L"
                                    ? "Laki-laki"
                                    : student.gender === "P"
                                      ? "Perempuan"
                                      : "-"}</span
                            >
                        {/if}
                    </div>
                    <div class="info-item">
                        <label>Agama</label>
                        {#if isEditing}
                            <select bind:value={editForm.religion}>
                                <option value="">Pilih...</option>
                                <option value="islam">Islam</option>
                                <option value="kristen">Kristen</option>
                                <option value="katolik">Katolik</option>
                                <option value="hindu">Hindu</option>
                                <option value="buddha">Buddha</option>
                                <option value="konghucu">Konghucu</option>
                            </select>
                        {:else}
                            <span class="capitalize"
                                >{student.religion || "-"}</span
                            >
                        {/if}
                    </div>
                    <div class="info-item full-width">
                        <label>Alamat</label>
                        {#if isEditing}
                            <textarea bind:value={editForm.address}></textarea>
                        {:else}
                            <span>{student.address || "-"}</span>
                        {/if}
                    </div>
                </div>
            {:else if activeTab === "parents"}
                <div class="info-grid">
                    <div class="info-item">
                        <label>Nama Ayah</label>
                        {#if isEditing}
                            <input
                                type="text"
                                bind:value={editForm.father_name}
                            />
                        {:else}
                            <span>{student.father_name || "-"}</span>
                        {/if}
                    </div>
                    <div class="info-item">
                        <label>Nama Ibu</label>
                        {#if isEditing}
                            <input
                                type="text"
                                bind:value={editForm.mother_name}
                            />
                        {:else}
                            <span>{student.mother_name || "-"}</span>
                        {/if}
                    </div>
                    <div class="info-item">
                        <label>Nama Wali</label>
                        {#if isEditing}
                            <input
                                type="text"
                                bind:value={editForm.guardian_name}
                            />
                        {:else}
                            <span>{student.guardian_name || "-"}</span>
                        {/if}
                    </div>
                    <div class="info-item">
                        <label>No. Telepon Orang Tua</label>
                        {#if isEditing}
                            <input
                                type="tel"
                                bind:value={editForm.parent_phone}
                            />
                        {:else}
                            <span>{student.parent_phone || "-"}</span>
                        {/if}
                    </div>
                </div>
            {:else if activeTab === "enrollment"}
                <div class="enrollment-section">
                    {#if student.enrollments && student.enrollments.length > 0}
                        <div class="enrollment-list">
                            {#each student.enrollments as enrollment}
                                <div
                                    class="enrollment-item"
                                    class:active={enrollment.status ===
                                        "active"}
                                >
                                    <div class="enrollment-class">
                                        <span class="class-name"
                                            >{enrollment.class?.name ||
                                                "-"}</span
                                        >
                                        <span class="academic-year"
                                            >{enrollment.academic_year}</span
                                        >
                                    </div>
                                    <div class="enrollment-status">
                                        <span
                                            class="status-badge {enrollment.status}"
                                        >
                                            {enrollment.status === "active"
                                                ? "Aktif"
                                                : enrollment.status ===
                                                    "graduated"
                                                  ? "Lulus"
                                                  : enrollment.status ===
                                                      "transferred"
                                                    ? "Pindah"
                                                    : "Keluar"}
                                        </span>
                                    </div>
                                    <div class="enrollment-dates">
                                        <span
                                            >Masuk: {formatDate(
                                                enrollment.enrolled_at,
                                            )}</span
                                        >
                                        {#if enrollment.left_at}
                                            <span
                                                >Keluar: {formatDate(
                                                    enrollment.left_at,
                                                )}</span
                                            >
                                        {/if}
                                    </div>
                                </div>
                            {/each}
                        </div>
                    {:else}
                        <div class="empty-enrollment">
                            <p>Belum terdaftar di kelas manapun</p>
                        </div>
                    {/if}
                </div>
            {/if}
        </div>
    {:else}
        <div class="not-found glass-card">
            <h3>Siswa Tidak Ditemukan</h3>
            <button class="btn-back-link" on:click={goBack}
                >← Kembali ke Daftar</button
            >
        </div>
    {/if}
</div>

<style>
    /* Container */
    .student-detail-container {
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
        margin-bottom: 1.5rem;
        background: #fee2e2;
        border-color: #fecaca;
        color: #dc2626;
    }

    .error-message button {
        background: none;
        border: none;
        color: #dc2626;
        font-size: 1.25rem;
        cursor: pointer;
    }

    /* Loading */
    .loading-state {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 4rem;
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

    .profile-info {
        flex: 1;
    }

    .profile-info h1 {
        color: #111827;
        font-size: 1.75rem;
        margin: 0 0 0.75rem;
    }

    .profile-badges {
        display: flex;
        flex-wrap: wrap;
        gap: 0.5rem;
    }

    .badge {
        display: inline-block;
        padding: 0.35rem 0.75rem;
        border-radius: 20px;
        font-size: 0.85rem;
        background: #f3f4f6;
        color: #374151;
    }

    .badge.nis {
        background: #dbeafe;
        color: #1d4ed8;
    }

    .badge.nisn {
        background: #d1fae5;
        color: #059669;
    }

    .badge.gender.male {
        background: #dbeafe;
        color: #1d4ed8;
    }

    .badge.gender.female {
        background: #fce7f3;
        color: #be185d;
    }

    .profile-actions {
        display: flex;
        gap: 0.75rem;
    }

    .btn-edit,
    .btn-cancel,
    .btn-save {
        padding: 0.75rem 1.25rem;
        border-radius: 10px;
        font-size: 0.95rem;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.2s;
    }

    .btn-edit {
        background: #00ced1;
        border: none;
        color: white;
    }

    .btn-edit:hover {
        background: #00b5b8;
    }

    .btn-cancel {
        background: #ffffff;
        border: 1px solid #e5e7eb;
        color: #374151;
    }

    .btn-cancel:hover {
        background: #f3f4f6;
    }

    .btn-save {
        background: #00ced1;
        border: none;
        color: white;
    }

    .btn-save:disabled {
        opacity: 0.6;
        cursor: not-allowed;
    }

    /* Tabs */
    .tabs {
        display: flex;
        padding: 0.5rem;
        margin-bottom: 1.5rem;
        gap: 0.5rem;
    }

    .tab {
        flex: 1;
        background: transparent;
        border: none;
        color: #6b7280;
        padding: 1rem;
        border-radius: 10px;
        cursor: pointer;
        font-size: 0.95rem;
        font-weight: 500;
        transition: all 0.2s;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
    }

    .tab:hover {
        background: #f3f4f6;
        color: #111827;
    }

    .tab.active {
        background: #00ced1;
        color: white;
    }

    /* Tab Content */
    .tab-content {
        padding: 2rem;
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

    .info-item.full-width {
        grid-column: 1 / -1;
    }

    .info-item label {
        color: #6b7280;
        font-size: 0.85rem;
        text-transform: uppercase;
        letter-spacing: 0.5px;
    }

    .info-item span {
        color: #111827;
        font-size: 1rem;
    }

    .info-item .capitalize {
        text-transform: capitalize;
    }

    .info-item input,
    .info-item select,
    .info-item textarea {
        padding: 0.75rem;
        background: #ffffff;
        border: 1px solid #e5e7eb;
        border-radius: 8px;
        color: #111827;
        font-size: 1rem;
        outline: none;
    }

    .info-item input:focus,
    .info-item select:focus,
    .info-item textarea:focus {
        border-color: #00ced1;
        box-shadow: 0 0 0 3px rgba(0, 206, 209, 0.1);
    }

    .info-item select option {
        background: #ffffff;
        color: #111827;
    }

    .info-item textarea {
        resize: vertical;
        min-height: 80px;
    }

    /* Enrollment */
    .enrollment-list {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    .enrollment-item {
        display: grid;
        grid-template-columns: 1fr auto auto;
        gap: 1rem;
        padding: 1.25rem;
        background: #f9fafb;
        border: 1px solid #e5e7eb;
        border-radius: 12px;
        align-items: center;
    }

    .enrollment-item.active {
        border-left: 4px solid #00ced1;
    }

    .enrollment-class {
        display: flex;
        flex-direction: column;
        gap: 0.25rem;
    }

    .class-name {
        color: #111827;
        font-weight: 600;
        font-size: 1.05rem;
    }

    .academic-year {
        color: #6b7280;
        font-size: 0.9rem;
    }

    .status-badge {
        padding: 0.35rem 0.75rem;
        border-radius: 20px;
        font-size: 0.85rem;
    }

    .status-badge.active {
        background: #d1fae5;
        color: #059669;
    }

    .status-badge.graduated {
        background: #dbeafe;
        color: #1d4ed8;
    }

    .status-badge.transferred {
        background: #fef3c7;
        color: #d97706;
    }

    .status-badge.dropped {
        background: #f3f4f6;
        color: #6b7280;
    }

    .enrollment-dates {
        display: flex;
        flex-direction: column;
        gap: 0.25rem;
        color: #6b7280;
        font-size: 0.85rem;
        text-align: right;
    }

    .empty-enrollment {
        text-align: center;
        padding: 3rem;
        color: #6b7280;
    }

    .empty-icon {
        font-size: 3rem;
        margin-bottom: 1rem;
    }

    /* Not Found */
    .not-found {
        text-align: center;
        padding: 4rem;
        color: #6b7280;
    }

    .not-found h3 {
        margin: 1rem 0;
        color: #111827;
    }

    .btn-back-link {
        background: #00ced1;
        border: none;
        color: white;
        padding: 0.75rem 1.25rem;
        border-radius: 10px;
        cursor: pointer;
        margin-top: 1rem;
    }

    .btn-back-link:hover {
        background: #00b5b8;
    }

    /* Responsive */
    @media (max-width: 768px) {
        .student-detail-container {
            padding: 1rem;
        }

        .profile-header {
            flex-direction: column;
            text-align: center;
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

        .enrollment-item {
            grid-template-columns: 1fr;
            gap: 0.75rem;
        }

        .enrollment-dates {
            text-align: left;
        }

        .tabs {
            flex-direction: column;
        }
    }
</style>
