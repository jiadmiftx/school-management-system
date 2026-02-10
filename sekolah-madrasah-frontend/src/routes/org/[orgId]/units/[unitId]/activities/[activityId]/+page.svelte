<script lang="ts">
    import { page } from "$app/stores";
    import { api, PageHeader, ActionButton, Modal } from "$lib";
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";

    // Get IDs from URL
    $: orgId = $page.params.orgId;
    $: unitId = $page.params.unitId;
    $: activityId = $page.params.activityId;

    // State
    let activity: any = null;
    let teachers: any[] = [];
    let students: any[] = [];
    let allTeachers: any[] = [];
    let allStudents: any[] = [];
    let isLoading = true;
    let error = "";

    // Modals
    let showAssignTeacherModal = false;
    let showEnrollStudentModal = false;
    let showRemoveTeacherModal = false;
    let showRemoveStudentModal = false;
    let selectedTeacher: any = null;
    let selectedStudent: any = null;
    let selectedTeacherId = "";
    let selectedStudentId = "";
    let teacherRole = "pembina";

    // Role check
    let canManage = true;

    // Activity types
    const activityTypes: Record<string, string> = {
        ekstrakurikuler: "Ekstrakurikuler",
        kajian: "Kajian",
        event: "Event",
    };

    // Categories
    const categories: Record<string, string> = {
        halaqah: "Halaqah",
        tahsin: "Tahsin",
        daurah: "Daurah",
        olahraga: "Olahraga",
        seni: "Seni",
        akademik: "Akademik",
    };

    // Days of week
    const daysOfWeek = [
        "Minggu",
        "Senin",
        "Selasa",
        "Rabu",
        "Kamis",
        "Jumat",
        "Sabtu",
    ];

    async function checkUserRole() {
        try {
            const membershipRes: any = await api.get("/users/me/memberships");
            const memberships = membershipRes.data?.unit_memberships || [];
            const membership = memberships.find(
                (m: any) => m.unit_id === unitId,
            );
            if (membership) {
                canManage = ["owner", "admin", "staff"].includes(
                    membership.role,
                );
            }
        } catch (err) {
            console.error("Failed to check role:", err);
        }
    }

    async function loadActivity() {
        try {
            isLoading = true;
            error = "";
            const response: any = await api.get(`/activities/${activityId}`);
            activity = response.data;
        } catch (err: any) {
            error = err.message || "Gagal memuat data kegiatan";
        } finally {
            isLoading = false;
        }
    }

    async function loadTeachers() {
        try {
            const response: any = await api.get(
                `/activities/${activityId}/teachers`,
            );
            teachers = response.data || [];
        } catch (err: any) {
            console.error("Failed to load teachers:", err);
        }
    }

    async function loadStudents() {
        try {
            const response: any = await api.get(
                `/activities/${activityId}/students`,
            );
            students = response.data || [];
        } catch (err: any) {
            console.error("Failed to load students:", err);
        }
    }

    async function loadAllTeachers() {
        try {
            const response: any = await api.get(
                `/units/${unitId}/teachers?limit=100`,
            );
            allTeachers = response.data || [];
        } catch (err: any) {
            console.error("Failed to load all teachers:", err);
        }
    }

    async function loadAllStudents() {
        try {
            const response: any = await api.get(
                `/units/${unitId}/students?limit=100`,
            );
            allStudents = response.data || [];
        } catch (err: any) {
            console.error("Failed to load all students:", err);
        }
    }

    async function assignTeacher() {
        if (!selectedTeacherId) return;
        try {
            await api.post(`/activities/${activityId}/teachers`, {
                teacher_profile_id: selectedTeacherId,
                role: teacherRole,
            });
            showAssignTeacherModal = false;
            selectedTeacherId = "";
            teacherRole = "pembina";
            loadTeachers();
        } catch (err: any) {
            error = err.message || "Gagal menambahkan pembina";
        }
    }

    async function removeTeacher() {
        if (!selectedTeacher) return;
        try {
            await api.delete(
                `/activities/${activityId}/teachers/${selectedTeacher.teacher_profile_id}`,
            );
            showRemoveTeacherModal = false;
            selectedTeacher = null;
            loadTeachers();
        } catch (err: any) {
            error = err.message || "Gagal menghapus pembina";
        }
    }

    async function enrollStudent() {
        if (!selectedStudentId) return;
        try {
            await api.post(`/activities/${activityId}/students`, {
                student_profile_id: selectedStudentId,
                is_mandatory: false,
            });
            showEnrollStudentModal = false;
            selectedStudentId = "";
            loadStudents();
        } catch (err: any) {
            error = err.message || "Gagal mendaftarkan siswa";
        }
    }

    async function removeStudent() {
        if (!selectedStudent) return;
        try {
            await api.delete(
                `/activities/${activityId}/students/${selectedStudent.student_profile_id}`,
            );
            showRemoveStudentModal = false;
            selectedStudent = null;
            loadStudents();
        } catch (err: any) {
            error = err.message || "Gagal menghapus siswa";
        }
    }

    function openAssignTeacherModal() {
        loadAllTeachers();
        showAssignTeacherModal = true;
    }

    function openEnrollStudentModal() {
        loadAllStudents();
        showEnrollStudentModal = true;
    }

    function getDayLabel(day: number | null) {
        if (day === null || day === undefined) return "-";
        return daysOfWeek[day] || "-";
    }

    function getDaysLabel(days: number[] | null) {
        if (!days || days.length === 0) return "-";
        return days.map((d) => daysOfWeek[d] || d).join(", ");
    }

    function formatDate(dateStr: string | null) {
        if (!dateStr) return "-";
        const date = new Date(dateStr);
        return date.toLocaleDateString("id-ID", {
            day: "numeric",
            month: "long",
            year: "numeric",
        });
    }

    function formatCurrency(amount: number | null) {
        if (!amount) return "Gratis";
        return new Intl.NumberFormat("id-ID", {
            style: "currency",
            currency: "IDR",
            minimumFractionDigits: 0,
        }).format(amount);
    }

    onMount(async () => {
        await checkUserRole();
        await loadActivity();
        await Promise.all([loadTeachers(), loadStudents()]);
    });
</script>

<svelte:head>
    <title>{activity?.name || "Detail Kegiatan"} - Sekolah Madrasah</title>
</svelte:head>

<div class="page-container">
    <!-- Back Button -->
    <a href="/org/{orgId}/units/{unitId}/activities" class="back-btn">
        <svg
            xmlns="http://www.w3.org/2000/svg"
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
        >
            <path d="M19 12H5M12 19l-7-7 7-7" />
        </svg>
        Kembali
    </a>

    {#if isLoading}
        <div class="loading-state glass-card">
            <div class="loader"></div>
            <p>Memuat data kegiatan...</p>
        </div>
    {:else if error}
        <div class="error-message glass-card">
            <span>⚠️ {error}</span>
            <button on:click={() => (error = "")}>×</button>
        </div>
    {:else if activity}
        <!-- Activity Header -->
        <div class="activity-header glass-card">
            <div class="header-content">
                <div class="header-icon">
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="32"
                        height="32"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="white"
                        stroke-width="1.5"
                    >
                        <path
                            d="M4 15s1-1 4-1 5 2 8 2 4-1 4-1V3s-1 1-4 1-5-2-8-2-4 1-4 1z"
                        />
                        <path d="M4 22v-7" />
                    </svg>
                </div>
                <div class="header-text">
                    <h1>{activity.name}</h1>
                    <div class="header-meta">
                        <span class="type-badge type-{activity.type}">
                            {activityTypes[activity.type] || activity.type}
                        </span>
                        <span
                            class="status-badge"
                            class:active={activity.is_active}
                        >
                            {activity.is_active ? "Aktif" : "Nonaktif"}
                        </span>
                    </div>
                </div>
            </div>
        </div>

        <!-- Activity Info -->
        <div class="info-grid">
            <div class="info-card glass-card">
                <h3>📅 Jadwal</h3>
                <div class="info-content">
                    {#if activity.start_date || activity.end_date}
                        <p>
                            <strong>Periode:</strong>
                            {formatDate(activity.start_date)} - {activity.end_date
                                ? formatDate(activity.end_date)
                                : "Ongoing"}
                        </p>
                    {/if}
                    {#if activity.recurrence_type === "weekly"}
                        <p>
                            <strong>Setiap:</strong>
                            {getDaysLabel(activity.recurrence_days)}
                        </p>
                    {:else if activity.recurrence_type === "daily"}
                        <p><strong>Frekuensi:</strong> Setiap Hari</p>
                    {:else if activity.recurrence_type === "monthly"}
                        <p>
                            <strong>Tanggal:</strong>
                            {activity.recurrence_days?.join(", ") || "-"}
                        </p>
                    {:else}
                        <p><strong>Jadwal:</strong> Tidak Berulang</p>
                    {/if}
                    <p>
                        <strong>Waktu:</strong>
                        {activity.start_time || "-"} - {activity.end_time ||
                            "-"}
                    </p>
                </div>
            </div>
            <div class="info-card glass-card">
                <h3>📍 Lokasi & Kapasitas</h3>
                <div class="info-content">
                    <p>
                        <strong>Lokasi:</strong>
                        {activity.location || "Belum ditentukan"}
                    </p>
                    <p>
                        <strong>Maks. Peserta:</strong>
                        {activity.max_participants || "Tidak terbatas"}
                    </p>
                    <p>
                        <strong>Biaya:</strong>
                        {formatCurrency(activity.fee)}
                    </p>
                    {#if activity.category}
                        <p>
                            <strong>Kategori:</strong>
                            {categories[activity.category] || activity.category}
                        </p>
                    {/if}
                </div>
            </div>
            <div class="info-card glass-card full-width">
                <h3>📝 Deskripsi</h3>
                <p class="description-text">
                    {activity.description || "Tidak ada deskripsi"}
                </p>
            </div>
        </div>

        <!-- Teachers Section -->
        <div class="section glass-card">
            <div class="section-header">
                <h2>👨‍🏫 Pembina ({teachers.length})</h2>
                {#if canManage}
                    <button
                        class="btn-secondary"
                        on:click={openAssignTeacherModal}
                    >
                        + Tambah Pembina
                    </button>
                {/if}
            </div>
            {#if teachers.length === 0}
                <p class="empty-text">Belum ada pembina yang ditugaskan</p>
            {:else}
                <div class="member-list">
                    {#each teachers as teacher}
                        <div class="member-card">
                            <div class="avatar">
                                {teacher.teacher_profile?.user?.full_name?.charAt(
                                    0,
                                ) || "T"}
                            </div>
                            <div class="member-info">
                                <span class="member-name">
                                    {teacher.teacher_profile?.user?.full_name ||
                                        "-"}
                                </span>
                                <span class="member-role"
                                    >{teacher.role || "Pembina"}</span
                                >
                            </div>
                            {#if canManage}
                                <button
                                    class="btn-remove"
                                    on:click={() => {
                                        selectedTeacher = teacher;
                                        showRemoveTeacherModal = true;
                                    }}
                                >
                                    ×
                                </button>
                            {/if}
                        </div>
                    {/each}
                </div>
            {/if}
        </div>

        <!-- Students Section -->
        <div class="section glass-card">
            <div class="section-header">
                <h2>👨‍🎓 Peserta ({students.length})</h2>
                {#if canManage}
                    <button
                        class="btn-secondary"
                        on:click={openEnrollStudentModal}
                    >
                        + Daftarkan Siswa
                    </button>
                {/if}
            </div>
            {#if students.length === 0}
                <p class="empty-text">Belum ada siswa yang terdaftar</p>
            {:else}
                <div class="member-list">
                    {#each students as student}
                        <div class="member-card">
                            <div class="avatar student">
                                {student.student_profile?.user?.full_name?.charAt(
                                    0,
                                ) || "S"}
                            </div>
                            <div class="member-info">
                                <span class="member-name">
                                    {student.student_profile?.user?.full_name ||
                                        "-"}
                                </span>
                                <span class="member-role">
                                    NIS: {student.student_profile?.nis || "-"}
                                </span>
                            </div>
                            {#if canManage}
                                <button
                                    class="btn-remove"
                                    on:click={() => {
                                        selectedStudent = student;
                                        showRemoveStudentModal = true;
                                    }}
                                >
                                    ×
                                </button>
                            {/if}
                        </div>
                    {/each}
                </div>
            {/if}
        </div>
    {/if}
</div>

<!-- Assign Teacher Modal -->
<Modal
    show={showAssignTeacherModal}
    title="Tambah Pembina"
    size="default"
    on:close={() => (showAssignTeacherModal = false)}
>
    <div class="form-grid">
        <div class="form-group full-width">
            <label for="teacher">Pilih Guru</label>
            <select id="teacher" bind:value={selectedTeacherId}>
                <option value="">-- Pilih Guru --</option>
                {#each allTeachers as teacher}
                    <option value={teacher.id}>
                        {teacher.user?.full_name || "-"} ({teacher.nip || "-"})
                    </option>
                {/each}
            </select>
        </div>
        <div class="form-group full-width">
            <label for="role">Peran</label>
            <select id="role" bind:value={teacherRole}>
                <option value="pembina">Pembina</option>
                <option value="pengisi">Pengisi Materi</option>
                <option value="koordinator">Koordinator</option>
            </select>
        </div>
    </div>
    <svelte:fragment slot="actions">
        <button
            class="btn-secondary"
            on:click={() => (showAssignTeacherModal = false)}>Batal</button
        >
        <button
            class="btn-primary"
            on:click={assignTeacher}
            disabled={!selectedTeacherId}
        >
            Tambahkan
        </button>
    </svelte:fragment>
</Modal>

<!-- Enroll Student Modal -->
<Modal
    show={showEnrollStudentModal}
    title="Daftarkan Siswa"
    size="default"
    on:close={() => (showEnrollStudentModal = false)}
>
    <div class="form-group">
        <label for="student">Pilih Siswa</label>
        <select id="student" bind:value={selectedStudentId}>
            <option value="">-- Pilih Siswa --</option>
            {#each allStudents as student}
                <option value={student.id}>
                    {student.user?.full_name || "-"} ({student.nis || "-"})
                </option>
            {/each}
        </select>
    </div>
    <svelte:fragment slot="actions">
        <button
            class="btn-secondary"
            on:click={() => (showEnrollStudentModal = false)}>Batal</button
        >
        <button
            class="btn-primary"
            on:click={enrollStudent}
            disabled={!selectedStudentId}
        >
            Daftarkan
        </button>
    </svelte:fragment>
</Modal>

<!-- Remove Teacher Modal -->
<Modal
    show={showRemoveTeacherModal}
    title="Hapus Pembina"
    size="small"
    on:close={() => (showRemoveTeacherModal = false)}
>
    <div class="delete-confirm">
        <p>Apakah Anda yakin ingin menghapus pembina:</p>
        <p class="delete-name">
            {selectedTeacher?.teacher_profile?.user?.full_name || "-"}
        </p>
    </div>
    <svelte:fragment slot="actions">
        <button
            class="btn-secondary"
            on:click={() => (showRemoveTeacherModal = false)}>Batal</button
        >
        <button class="btn-danger" on:click={removeTeacher}>Hapus</button>
    </svelte:fragment>
</Modal>

<!-- Remove Student Modal -->
<Modal
    show={showRemoveStudentModal}
    title="Hapus Siswa"
    size="small"
    on:close={() => (showRemoveStudentModal = false)}
>
    <div class="delete-confirm">
        <p>Apakah Anda yakin ingin menghapus siswa:</p>
        <p class="delete-name">
            {selectedStudent?.student_profile?.user?.full_name || "-"}
        </p>
    </div>
    <svelte:fragment slot="actions">
        <button
            class="btn-secondary"
            on:click={() => (showRemoveStudentModal = false)}>Batal</button
        >
        <button class="btn-danger" on:click={removeStudent}>Hapus</button>
    </svelte:fragment>
</Modal>

<style>
    .page-container {
        padding: 2rem;
        min-height: 100vh;
        background: #f5f5f7;
    }

    .back-btn {
        display: inline-flex;
        align-items: center;
        gap: 0.5rem;
        margin-bottom: 1.5rem;
        padding: 0.5rem 1rem;
        color: #374151;
        text-decoration: none;
        font-weight: 500;
        background: #ffffff;
        border: 1px solid #e5e7eb;
        border-radius: 8px;
        transition: all 0.2s;
    }

    .back-btn:hover {
        background: #f3f4f6;
        border-color: #d1d5db;
    }

    .glass-card {
        background: #ffffff;
        border: 1px solid #e5e7eb;
        border-radius: 16px;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
    }

    /* Loading & Error */
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

    /* Activity Header */
    .activity-header {
        padding: 1.5rem 2rem;
        margin-bottom: 1.5rem;
    }

    .header-content {
        display: flex;
        align-items: center;
        gap: 1rem;
    }

    .header-icon {
        width: 56px;
        height: 56px;
        background: #00ced1;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .header-text h1 {
        font-size: 1.75rem;
        font-weight: 700;
        color: #111827;
        margin: 0 0 0.5rem;
    }

    .header-meta {
        display: flex;
        gap: 0.5rem;
    }

    .type-badge {
        display: inline-block;
        padding: 0.25rem 0.75rem;
        border-radius: 9999px;
        font-size: 0.75rem;
        font-weight: 500;
    }

    .type-ekstrakurikuler {
        background: #dbeafe;
        color: #3b82f6;
    }
    .type-kajian {
        background: #dcfce7;
        color: #16a34a;
    }
    .type-event {
        background: #fef3c7;
        color: #d97706;
    }

    .status-badge {
        display: inline-block;
        padding: 0.25rem 0.75rem;
        border-radius: 9999px;
        font-size: 0.75rem;
        font-weight: 500;
        background: #f3f4f6;
        color: #6b7280;
    }

    .status-badge.active {
        background: #dcfce7;
        color: #16a34a;
    }

    /* Info Grid */
    .info-grid {
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        gap: 1.5rem;
        margin-bottom: 1.5rem;
    }

    .info-card {
        padding: 1.5rem;
    }

    .info-card.full-width {
        grid-column: span 2;
    }

    .info-card h3 {
        font-size: 1rem;
        font-weight: 600;
        color: #374151;
        margin: 0 0 1rem;
    }

    .info-content p {
        margin: 0.5rem 0;
        color: #6b7280;
    }

    .description-text {
        color: #6b7280;
        line-height: 1.6;
    }

    /* Section */
    .section {
        padding: 1.5rem;
        margin-bottom: 1.5rem;
    }

    .section-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1rem;
    }

    .section-header h2 {
        font-size: 1.125rem;
        font-weight: 600;
        color: #374151;
        margin: 0;
    }

    .empty-text {
        color: #9ca3af;
        text-align: center;
        padding: 2rem;
    }

    /* Member List */
    .member-list {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
        gap: 1rem;
    }

    .member-card {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        padding: 0.75rem;
        background: #f9fafb;
        border-radius: 10px;
    }

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

    .avatar.student {
        background: #8b5cf6;
    }

    .member-info {
        flex: 1;
        display: flex;
        flex-direction: column;
    }

    .member-name {
        font-weight: 500;
        color: #111827;
    }

    .member-role {
        font-size: 0.75rem;
        color: #9ca3af;
    }

    .btn-remove {
        width: 28px;
        height: 28px;
        border: none;
        background: #fee2e2;
        color: #dc2626;
        border-radius: 6px;
        cursor: pointer;
        font-size: 1rem;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .btn-remove:hover {
        background: #fecaca;
    }

    /* Buttons */
    .btn-primary {
        background: #00ced1;
        color: white;
        border: none;
        padding: 0.75rem 1.5rem;
        border-radius: 12px;
        font-weight: 600;
        cursor: pointer;
    }

    .btn-primary:hover:not(:disabled) {
        background: #00b5b8;
    }

    .btn-primary:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }

    .btn-secondary {
        background: #f3f4f6;
        color: #374151;
        border: none;
        padding: 0.5rem 1rem;
        border-radius: 8px;
        font-weight: 500;
        cursor: pointer;
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
    }

    .btn-danger:hover {
        background: #b91c1c;
    }

    /* Form */
    .form-grid {
        display: grid;
        grid-template-columns: 1fr;
        gap: 1rem;
    }

    .form-group {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .form-group.full-width {
        grid-column: span 1;
    }

    .form-group label {
        font-weight: 500;
        color: #374151;
        font-size: 0.875rem;
    }

    .form-group select {
        padding: 0.75rem;
        border: 1px solid #e5e7eb;
        border-radius: 8px;
        font-size: 1rem;
        background: #ffffff;
    }

    .form-group select:focus {
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

    @media (max-width: 768px) {
        .page-container {
            padding: 1rem;
        }

        .info-grid {
            grid-template-columns: 1fr;
        }

        .member-list {
            grid-template-columns: 1fr;
        }
    }
</style>
