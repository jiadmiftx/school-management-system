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
    let activities: any[] = [];
    let isLoading = true;
    let error = "";
    let searchQuery = "";
    let filterType = "";

    // Modals
    let showCreateModal = false;
    let showEditModal = false;
    let showDeleteModal = false;
    let selectedActivity: any = null;

    // Pagination
    let currentPage = 1;
    let totalPages = 1;
    let totalItems = 0;
    let pageSize = 10;

    // Create/Edit form - Updated with new fields
    let formData = {
        name: "",
        type: "ekstrakurikuler",
        category: "",
        description: "",
        start_date: "",
        end_date: "",
        recurrence_type: "none",
        recurrence_days: [] as number[],
        start_time: "",
        end_time: "",
        location: "",
        max_participants: null as number | null,
        fee: null as number | null,
    };

    // Activity types
    const activityTypes = [
        { value: "ekstrakurikuler", label: "Ekstrakurikuler" },
        { value: "kajian", label: "Kajian" },
        { value: "event", label: "Event" },
    ];

    // Categories
    const categories = [
        { value: "", label: "-- Pilih Kategori --" },
        { value: "halaqah", label: "Halaqah" },
        { value: "tahsin", label: "Tahsin" },
        { value: "daurah", label: "Daurah" },
        { value: "olahraga", label: "Olahraga" },
        { value: "seni", label: "Seni" },
        { value: "akademik", label: "Akademik" },
    ];

    // Days of week
    const daysOfWeek = [
        { value: 0, label: "Minggu" },
        { value: 1, label: "Senin" },
        { value: 2, label: "Selasa" },
        { value: 3, label: "Rabu" },
        { value: 4, label: "Kamis" },
        { value: 5, label: "Jumat" },
        { value: 6, label: "Sabtu" },
    ];

    // Recurrence types
    const recurrenceTypes = [
        { value: "none", label: "Tidak Berulang" },
        { value: "daily", label: "Harian" },
        { value: "weekly", label: "Mingguan" },
        { value: "monthly", label: "Bulanan" },
    ];

    // Role check
    let canManage = true;

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

    async function loadActivities() {
        try {
            isLoading = true;
            error = "";
            let url = `/units/${unitId}/activities?page=${currentPage}&limit=${pageSize}`;
            if (filterType) url += `&type=${filterType}`;
            if (searchQuery) url += `&search=${searchQuery}`;

            const response: any = await api.get(url);

            // Handle nested response structure
            const data = response.data?.data || response.data || [];
            activities = Array.isArray(data) ? data : [];

            const total = response.data?.total || activities.length;
            totalItems = total;
            totalPages = Math.ceil(total / pageSize) || 1;
        } catch (err: any) {
            error = err.message || "Gagal memuat data kegiatan";
            activities = [];
        } finally {
            isLoading = false;
        }
    }

    async function createActivity() {
        try {
            const payload: any = {
                name: formData.name,
                type: formData.type,
                category: formData.category || undefined,
                description: formData.description || undefined,
                start_date: formData.start_date || undefined,
                end_date: formData.end_date || undefined,
                recurrence_type: formData.recurrence_type,
                recurrence_days:
                    formData.recurrence_days.length > 0
                        ? formData.recurrence_days
                        : undefined,
                start_time: formData.start_time || undefined,
                end_time: formData.end_time || undefined,
                location: formData.location || undefined,
                max_participants: formData.max_participants || undefined,
                fee: formData.fee || undefined,
            };

            await api.post(`/units/${unitId}/activities`, payload);
            showCreateModal = false;
            resetForm();
            loadActivities();
        } catch (err: any) {
            error = err.message || "Gagal menambah kegiatan";
        }
    }

    async function updateActivity() {
        if (!selectedActivity) return;
        try {
            const payload: any = {
                name: formData.name,
                type: formData.type,
                category: formData.category || undefined,
                description: formData.description || undefined,
                start_date: formData.start_date || undefined,
                end_date: formData.end_date || undefined,
                recurrence_type: formData.recurrence_type,
                recurrence_days:
                    formData.recurrence_days.length > 0
                        ? formData.recurrence_days
                        : undefined,
                start_time: formData.start_time || undefined,
                end_time: formData.end_time || undefined,
                location: formData.location || undefined,
                max_participants: formData.max_participants || undefined,
                fee: formData.fee || undefined,
            };

            await api.put(`/activities/${selectedActivity.id}`, payload);
            showEditModal = false;
            selectedActivity = null;
            resetForm();
            loadActivities();
        } catch (err: any) {
            error = err.message || "Gagal mengupdate kegiatan";
        }
    }

    async function deleteActivity() {
        if (!selectedActivity) return;
        try {
            await api.delete(`/activities/${selectedActivity.id}`);
            showDeleteModal = false;
            selectedActivity = null;
            loadActivities();
        } catch (err: any) {
            error = err.message || "Gagal menghapus kegiatan";
        }
    }

    function resetForm() {
        formData = {
            name: "",
            type: "ekstrakurikuler",
            category: "",
            description: "",
            start_date: "",
            end_date: "",
            recurrence_type: "none",
            recurrence_days: [],
            start_time: "",
            end_time: "",
            location: "",
            max_participants: null,
            fee: null,
        };
    }

    function openEditModal(activity: any) {
        selectedActivity = activity;
        formData = {
            name: activity.name || "",
            type: activity.type || "ekstrakurikuler",
            category: activity.category || "",
            description: activity.description || "",
            start_date: activity.start_date
                ? activity.start_date.split("T")[0]
                : "",
            end_date: activity.end_date ? activity.end_date.split("T")[0] : "",
            recurrence_type: activity.recurrence_type || "none",
            recurrence_days: activity.recurrence_days || [],
            start_time: activity.start_time || "",
            end_time: activity.end_time || "",
            location: activity.location || "",
            max_participants: activity.max_participants || null,
            fee: activity.fee || null,
        };
        showEditModal = true;
    }

    function openDeleteModal(activity: any) {
        selectedActivity = activity;
        showDeleteModal = true;
    }

    function handleSearch() {
        currentPage = 1;
        loadActivities();
    }

    function handlePageChange(e: CustomEvent) {
        currentPage = e.detail.page;
        loadActivities();
    }

    function getTypeLabel(type: string) {
        return activityTypes.find((t) => t.value === type)?.label || type;
    }

    function getTypeBadgeClass(type: string) {
        switch (type) {
            case "ekstrakurikuler":
                return "type-ekskul";
            case "kajian":
                return "type-kajian";
            case "event":
                return "type-event";
            default:
                return "";
        }
    }

    function getDaysLabel(days: number[] | null) {
        if (!days || days.length === 0) return "-";
        return days
            .map((d) => daysOfWeek.find((day) => day.value === d)?.label || d)
            .join(", ");
    }

    function formatTime(time: string | null) {
        if (!time) return "-";
        return time;
    }

    function toggleDay(day: number) {
        if (formData.recurrence_days.includes(day)) {
            formData.recurrence_days = formData.recurrence_days.filter(
                (d) => d !== day,
            );
        } else {
            formData.recurrence_days = [...formData.recurrence_days, day].sort(
                (a, b) => a - b,
            );
        }
    }

    onMount(async () => {
        await checkUserRole();
        loadActivities();
    });
</script>

<svelte:head>
    <title>Kegiatan - Sekolah Madrasah</title>
</svelte:head>

<div class="page-container">
    <!-- Page Header -->
    <PageHeader
        icon="activity"
        title="Kegiatan"
        subtitle="Kelola ekstrakurikuler, kajian, dan event"
        showAction={canManage}
        actionLabel="Tambah Kegiatan"
        on:action={() => (showCreateModal = true)}
    />

    <!-- Error Message -->
    {#if error}
        <div class="error-message glass-card">
            <span>⚠️ {error}</span>
            <button on:click={() => (error = "")}>×</button>
        </div>
    {/if}

    <!-- Search & Filter -->
    <div class="filter-bar glass-card">
        <div class="search-input-wrapper">
            <span class="search-icon">
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="18"
                    height="18"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                >
                    <circle cx="11" cy="11" r="8" /><path
                        d="m21 21-4.35-4.35"
                    />
                </svg>
            </span>
            <input
                type="text"
                placeholder="Cari nama kegiatan..."
                bind:value={searchQuery}
                on:keypress={(e) => e.key === "Enter" && handleSearch()}
            />
        </div>
        <select bind:value={filterType} on:change={handleSearch}>
            <option value="">Semua Jenis</option>
            {#each activityTypes as type}
                <option value={type.value}>{type.label}</option>
            {/each}
        </select>
        <button class="btn-search" on:click={handleSearch}>Cari</button>
        <div class="stats">
            <strong>{totalItems}</strong> kegiatan terdaftar
        </div>
    </div>

    <!-- Activities Table -->
    <div class="table-container glass-card">
        {#if isLoading}
            <div class="loading-state">
                <div class="loader"></div>
                <p>Memuat data kegiatan...</p>
            </div>
        {:else if activities.length === 0}
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
                        <path
                            d="M4 15s1-1 4-1 5 2 8 2 4-1 4-1V3s-1 1-4 1-5-2-8-2-4 1-4 1z"
                        />
                        <path d="M4 22v-7" />
                    </svg>
                </div>
                <h3>Belum Ada Kegiatan</h3>
                <p>Mulai dengan menambahkan kegiatan pertama</p>
                {#if canManage}
                    <button
                        class="btn-primary"
                        on:click={() => (showCreateModal = true)}
                    >
                        + Tambah Kegiatan
                    </button>
                {/if}
            </div>
        {:else}
            <table class="data-table">
                <thead>
                    <tr>
                        <th>Nama Kegiatan</th>
                        <th>Jenis</th>
                        <th>Jadwal</th>
                        <th>Waktu</th>
                        <th>Status</th>
                        <th>Aksi</th>
                    </tr>
                </thead>
                <tbody>
                    {#each activities as activity}
                        <tr>
                            <td>
                                <div class="activity-name">
                                    <a
                                        href="/org/{orgId}/units/{unitId}/activities/{activity.id}"
                                    >
                                        {activity.name}
                                    </a>
                                    {#if activity.description}
                                        <span class="description"
                                            >{activity.description}</span
                                        >
                                    {/if}
                                </div>
                            </td>
                            <td>
                                <span
                                    class="type-badge {getTypeBadgeClass(
                                        activity.type,
                                    )}"
                                >
                                    {getTypeLabel(activity.type)}
                                </span>
                            </td>
                            <td>
                                {#if activity.recurrence_type === "weekly"}
                                    Setiap {getDaysLabel(
                                        activity.recurrence_days,
                                    )}
                                {:else if activity.recurrence_type === "daily"}
                                    Setiap Hari
                                {:else if activity.recurrence_type === "monthly"}
                                    Tanggal {activity.recurrence_days?.join(
                                        ", ",
                                    ) || "-"}
                                {:else}
                                    Tidak Berulang
                                {/if}
                            </td>
                            <td>
                                {formatTime(activity.start_time)} - {formatTime(
                                    activity.end_time,
                                )}
                            </td>
                            <td>
                                <span
                                    class="status-badge"
                                    class:active={activity.is_active}
                                    class:inactive={!activity.is_active}
                                >
                                    {activity.is_active ? "Aktif" : "Nonaktif"}
                                </span>
                            </td>
                            <td>
                                <div class="actions">
                                    <ActionButton
                                        type="view"
                                        href="/org/{orgId}/units/{unitId}/activities/{activity.id}"
                                        title="Lihat Detail"
                                    />
                                    {#if canManage}
                                        <ActionButton
                                            type="edit"
                                            title="Edit"
                                            on:click={() =>
                                                openEditModal(activity)}
                                        />
                                        <ActionButton
                                            type="delete"
                                            title="Hapus"
                                            on:click={() =>
                                                openDeleteModal(activity)}
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
    title="Tambah Kegiatan Baru"
    size="large"
    on:close={() => (showCreateModal = false)}
>
    <form on:submit|preventDefault={createActivity}>
        <div class="form-grid">
            <div class="form-group full-width">
                <label for="name"
                    >Nama Kegiatan <span class="required">*</span></label
                >
                <input
                    type="text"
                    id="name"
                    bind:value={formData.name}
                    placeholder="Contoh: Pramuka, Futsal, Tahfidz"
                    required
                />
            </div>
            <div class="form-group">
                <label for="type"
                    >Jenis Kegiatan <span class="required">*</span></label
                >
                <select id="type" bind:value={formData.type} required>
                    {#each activityTypes as type}
                        <option value={type.value}>{type.label}</option>
                    {/each}
                </select>
            </div>
            <div class="form-group">
                <label for="category">Kategori</label>
                <select id="category" bind:value={formData.category}>
                    {#each categories as cat}
                        <option value={cat.value}>{cat.label}</option>
                    {/each}
                </select>
            </div>

            <!-- Date Range -->
            <div class="form-group">
                <label for="start_date">Tanggal Mulai</label>
                <input
                    type="date"
                    id="start_date"
                    bind:value={formData.start_date}
                />
            </div>
            <div class="form-group">
                <label for="end_date">Tanggal Berakhir</label>
                <input
                    type="date"
                    id="end_date"
                    bind:value={formData.end_date}
                />
            </div>

            <!-- Recurrence -->
            <div class="form-group">
                <label for="recurrence">Pengulangan</label>
                <select id="recurrence" bind:value={formData.recurrence_type}>
                    {#each recurrenceTypes as rec}
                        <option value={rec.value}>{rec.label}</option>
                    {/each}
                </select>
            </div>

            {#if formData.recurrence_type === "weekly"}
                <div class="form-group full-width">
                    <label>Pilih Hari (bisa lebih dari satu)</label>
                    <div class="day-selector">
                        {#each daysOfWeek as day}
                            <button
                                type="button"
                                class="day-btn"
                                class:selected={formData.recurrence_days.includes(
                                    day.value,
                                )}
                                on:click={() => toggleDay(day.value)}
                            >
                                {day.label}
                            </button>
                        {/each}
                    </div>
                </div>
            {/if}

            {#if formData.recurrence_type === "monthly"}
                <div class="form-group full-width">
                    <label for="monthly_days"
                        >Tanggal Bulanan (pisahkan dengan koma)</label
                    >
                    <input
                        type="text"
                        id="monthly_days"
                        placeholder="Contoh: 1, 15"
                        on:change={(e) => {
                            const input = e.currentTarget.value;
                            formData.recurrence_days = input
                                .split(",")
                                .map((s) => parseInt(s.trim()))
                                .filter((n) => !isNaN(n) && n >= 1 && n <= 31);
                        }}
                    />
                </div>
            {/if}

            <!-- Time -->
            <div class="form-group">
                <label for="start_time">Jam Mulai</label>
                <input
                    type="time"
                    id="start_time"
                    bind:value={formData.start_time}
                />
            </div>
            <div class="form-group">
                <label for="end_time">Jam Selesai</label>
                <input
                    type="time"
                    id="end_time"
                    bind:value={formData.end_time}
                />
            </div>

            <!-- Additional Info -->
            <div class="form-group full-width">
                <label for="location">Lokasi</label>
                <input
                    type="text"
                    id="location"
                    bind:value={formData.location}
                    placeholder="Contoh: Lapangan, Aula, Masjid"
                />
            </div>
            <div class="form-group">
                <label for="max_participants">Maks. Peserta</label>
                <input
                    type="number"
                    id="max_participants"
                    bind:value={formData.max_participants}
                    placeholder="Kosongkan jika tidak terbatas"
                    min="1"
                />
            </div>
            <div class="form-group">
                <label for="fee">Biaya (Rp)</label>
                <input
                    type="number"
                    id="fee"
                    bind:value={formData.fee}
                    placeholder="Kosongkan jika gratis"
                    min="0"
                    step="1000"
                />
            </div>

            <div class="form-group full-width">
                <label for="description">Deskripsi</label>
                <textarea
                    id="description"
                    bind:value={formData.description}
                    placeholder="Deskripsi kegiatan (opsional)"
                    rows="3"
                ></textarea>
            </div>
        </div>
    </form>
    <svelte:fragment slot="actions">
        <button class="btn-secondary" on:click={() => (showCreateModal = false)}
            >Batal</button
        >
        <button class="btn-primary" on:click={createActivity}>Simpan</button>
    </svelte:fragment>
</Modal>

<!-- Edit Modal -->
<Modal
    show={showEditModal}
    title="Edit Kegiatan"
    size="large"
    on:close={() => (showEditModal = false)}
>
    <form on:submit|preventDefault={updateActivity}>
        <div class="form-grid">
            <div class="form-group full-width">
                <label for="edit-name"
                    >Nama Kegiatan <span class="required">*</span></label
                >
                <input
                    type="text"
                    id="edit-name"
                    bind:value={formData.name}
                    required
                />
            </div>
            <div class="form-group">
                <label for="edit-type"
                    >Jenis Kegiatan <span class="required">*</span></label
                >
                <select id="edit-type" bind:value={formData.type} required>
                    {#each activityTypes as type}
                        <option value={type.value}>{type.label}</option>
                    {/each}
                </select>
            </div>
            <div class="form-group">
                <label for="edit-category">Kategori</label>
                <select id="edit-category" bind:value={formData.category}>
                    {#each categories as cat}
                        <option value={cat.value}>{cat.label}</option>
                    {/each}
                </select>
            </div>

            <!-- Date Range -->
            <div class="form-group">
                <label for="edit-start_date">Tanggal Mulai</label>
                <input
                    type="date"
                    id="edit-start_date"
                    bind:value={formData.start_date}
                />
            </div>
            <div class="form-group">
                <label for="edit-end_date">Tanggal Berakhir</label>
                <input
                    type="date"
                    id="edit-end_date"
                    bind:value={formData.end_date}
                />
            </div>

            <!-- Recurrence -->
            <div class="form-group">
                <label for="edit-recurrence">Pengulangan</label>
                <select
                    id="edit-recurrence"
                    bind:value={formData.recurrence_type}
                >
                    {#each recurrenceTypes as rec}
                        <option value={rec.value}>{rec.label}</option>
                    {/each}
                </select>
            </div>

            {#if formData.recurrence_type === "weekly"}
                <div class="form-group full-width">
                    <label>Pilih Hari (bisa lebih dari satu)</label>
                    <div class="day-selector">
                        {#each daysOfWeek as day}
                            <button
                                type="button"
                                class="day-btn"
                                class:selected={formData.recurrence_days.includes(
                                    day.value,
                                )}
                                on:click={() => toggleDay(day.value)}
                            >
                                {day.label}
                            </button>
                        {/each}
                    </div>
                </div>
            {/if}

            {#if formData.recurrence_type === "monthly"}
                <div class="form-group full-width">
                    <label for="edit-monthly_days"
                        >Tanggal Bulanan (pisahkan dengan koma)</label
                    >
                    <input
                        type="text"
                        id="edit-monthly_days"
                        value={formData.recurrence_days.join(", ")}
                        on:change={(e) => {
                            const input = e.currentTarget.value;
                            formData.recurrence_days = input
                                .split(",")
                                .map((s) => parseInt(s.trim()))
                                .filter((n) => !isNaN(n) && n >= 1 && n <= 31);
                        }}
                    />
                </div>
            {/if}

            <!-- Time -->
            <div class="form-group">
                <label for="edit-start">Jam Mulai</label>
                <input
                    type="time"
                    id="edit-start"
                    bind:value={formData.start_time}
                />
            </div>
            <div class="form-group">
                <label for="edit-end">Jam Selesai</label>
                <input
                    type="time"
                    id="edit-end"
                    bind:value={formData.end_time}
                />
            </div>

            <!-- Additional Info -->
            <div class="form-group full-width">
                <label for="edit-location">Lokasi</label>
                <input
                    type="text"
                    id="edit-location"
                    bind:value={formData.location}
                    placeholder="Contoh: Lapangan, Aula, Masjid"
                />
            </div>
            <div class="form-group">
                <label for="edit-max_participants">Maks. Peserta</label>
                <input
                    type="number"
                    id="edit-max_participants"
                    bind:value={formData.max_participants}
                    placeholder="Kosongkan jika tidak terbatas"
                    min="1"
                />
            </div>
            <div class="form-group">
                <label for="edit-fee">Biaya (Rp)</label>
                <input
                    type="number"
                    id="edit-fee"
                    bind:value={formData.fee}
                    placeholder="Kosongkan jika gratis"
                    min="0"
                    step="1000"
                />
            </div>

            <div class="form-group full-width">
                <label for="edit-desc">Deskripsi</label>
                <textarea
                    id="edit-desc"
                    bind:value={formData.description}
                    rows="3"
                ></textarea>
            </div>
        </div>
    </form>
    <svelte:fragment slot="actions">
        <button class="btn-secondary" on:click={() => (showEditModal = false)}
            >Batal</button
        >
        <button class="btn-primary" on:click={updateActivity}>Simpan</button>
    </svelte:fragment>
</Modal>

<!-- Delete Modal -->
<Modal
    show={showDeleteModal}
    title="Hapus Kegiatan"
    size="small"
    on:close={() => (showDeleteModal = false)}
>
    <div class="delete-confirm">
        <p>Apakah Anda yakin ingin menghapus kegiatan:</p>
        <p class="delete-name">{selectedActivity?.name || "-"}</p>
        <p class="delete-warning">Tindakan ini tidak dapat dibatalkan!</p>
    </div>
    <svelte:fragment slot="actions">
        <button class="btn-secondary" on:click={() => (showDeleteModal = false)}
            >Batal</button
        >
        <button class="btn-danger" on:click={deleteActivity}>Hapus</button>
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

    /* Filter Bar */
    .filter-bar {
        display: flex;
        align-items: center;
        gap: 1rem;
        padding: 1rem 1.5rem;
        margin-bottom: 1.5rem;
    }

    .search-input-wrapper {
        flex: 1;
        display: flex;
        align-items: center;
        background: #ffffff;
        border: 1px solid #e5e7eb;
        border-radius: 10px;
        padding: 0 1rem;
    }

    .search-icon {
        color: #9ca3af;
        display: flex;
        margin-right: 0.5rem;
    }

    .search-input-wrapper input {
        flex: 1;
        background: transparent;
        border: none;
        padding: 0.75rem 0;
        font-size: 1rem;
        outline: none;
    }

    .filter-bar select {
        padding: 0.75rem 1rem;
        border: 1px solid #e5e7eb;
        border-radius: 10px;
        background: white;
        font-size: 1rem;
        min-width: 160px;
    }

    .btn-search {
        background: #00ced1;
        border: none;
        color: white;
        padding: 0.75rem 1.5rem;
        border-radius: 10px;
        cursor: pointer;
        font-weight: 500;
    }

    .btn-search:hover {
        background: #00b5b8;
    }

    .stats {
        color: #6b7280;
        font-size: 0.875rem;
        white-space: nowrap;
    }

    .stats strong {
        color: #111827;
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

    /* Activity Name */
    .activity-name {
        display: flex;
        flex-direction: column;
        gap: 0.25rem;
    }

    .activity-name a {
        font-weight: 500;
        color: #00ced1;
        text-decoration: none;
    }

    .activity-name a:hover {
        text-decoration: underline;
    }

    .activity-name .description {
        font-size: 0.75rem;
        color: #9ca3af;
        max-width: 200px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    }

    /* Type Badge */
    .type-badge {
        display: inline-block;
        padding: 0.25rem 0.75rem;
        border-radius: 9999px;
        font-size: 0.75rem;
        font-weight: 500;
    }

    .type-ekskul {
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

    /* Status Badge */
    .status-badge {
        display: inline-block;
        padding: 0.25rem 0.75rem;
        border-radius: 9999px;
        font-size: 0.75rem;
        font-weight: 500;
    }

    .status-badge.active {
        background: #dcfce7;
        color: #16a34a;
    }

    .status-badge.inactive {
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

        .filter-bar {
            flex-wrap: wrap;
        }

        .form-grid {
            grid-template-columns: 1fr;
        }

        .form-group.full-width {
            grid-column: span 1;
        }
    }

    /* Day Selector */
    .day-selector {
        display: flex;
        flex-wrap: wrap;
        gap: 0.5rem;
    }

    .day-btn {
        padding: 0.5rem 1rem;
        border: 1px solid #e5e7eb;
        border-radius: 8px;
        background: #ffffff;
        color: #374151;
        cursor: pointer;
        transition: all 0.2s ease;
        font-size: 0.875rem;
    }

    .day-btn:hover {
        border-color: #00ced1;
        background: #f0fdfa;
    }

    .day-btn.selected {
        background: #00ced1;
        border-color: #00ced1;
        color: #ffffff;
    }
</style>
