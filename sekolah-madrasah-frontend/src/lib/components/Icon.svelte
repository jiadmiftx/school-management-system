<script lang="ts">
    export let name: string;
    export let size: number = 20;
    export let color: string = "currentColor"; // Alias for stroke
    export let stroke: string = "";
    export let fill: string = "none";
    export let strokeWidth: number = 1.5;

    // Use color as stroke if stroke not explicitly set
    $: effectiveStroke = stroke || color;

    // HugeIcons-inspired SVG paths
    const icons: Record<string, string> = {
        // View/Eye icon
        view: `<path d="M2.5 12C2.5 12 5.5 5 12 5C18.5 5 21.5 12 21.5 12C21.5 12 18.5 19 12 19C5.5 19 2.5 12 2.5 12Z" stroke-linecap="round" stroke-linejoin="round"/><circle cx="12" cy="12" r="3" stroke-linecap="round" stroke-linejoin="round"/>`,

        // Key icon
        key: `<path d="M21 2L19 4M19 4L22 7L18.5 10.5L15.5 7.5M19 4L15.5 7.5M11.39 11.61C12.3249 12.5353 12.8536 13.7974 12.8536 15.1136C12.8536 16.4299 12.3249 17.692 11.39 18.6173C10.4551 19.5427 9.18091 20.0659 7.85195 20.0659C6.52299 20.0659 5.24884 19.5427 4.31393 18.6173C3.37903 17.692 2.85034 16.4299 2.85034 15.1136C2.85034 13.7974 3.37903 12.5353 4.31393 11.61C5.24884 10.6846 6.52299 10.1614 7.85195 10.1614C9.18091 10.1614 10.4551 10.6846 11.39 11.61ZM11.39 11.61L15.5 7.5" stroke-linecap="round" stroke-linejoin="round"/>`,

        // Trash/Delete icon
        trash: `<path d="M4 7H20M10 11V17M14 11V17M5 7L6 19C6 19.5304 6.21071 20.0391 6.58579 20.4142C6.96086 20.7893 7.46957 21 8 21H16C16.5304 21 17.0391 20.7893 17.4142 20.4142C17.7893 20.0391 18 19.5304 18 19L19 7M9 7V4C9 3.73478 9.10536 3.48043 9.29289 3.29289C9.48043 3.10536 9.73478 3 10 3H14C14.2652 3 14.5196 3.10536 14.7071 3.29289C14.8946 3.48043 15 3.73478 15 4V7" stroke-linecap="round" stroke-linejoin="round"/>`,

        // Plus/Add icon
        plus: `<path d="M12 5V19M5 12H19" stroke-linecap="round" stroke-linejoin="round"/>`,

        // Search icon
        search: `<circle cx="11" cy="11" r="7" stroke-linecap="round" stroke-linejoin="round"/><path d="M21 21L16.5 16.5" stroke-linecap="round" stroke-linejoin="round"/>`,

        // User/Student icon
        user: `<circle cx="12" cy="8" r="4" stroke-linecap="round" stroke-linejoin="round"/><path d="M20 21C20 18.8783 19.1571 16.8434 17.6569 15.3431C16.1566 13.8429 14.1217 13 12 13C9.87827 13 7.84344 13.8429 6.34315 15.3431C4.84286 16.8434 4 18.8783 4 21" stroke-linecap="round" stroke-linejoin="round"/>`,

        // Users/Team icon
        users: `<circle cx="9" cy="7" r="3" stroke-linecap="round" stroke-linejoin="round"/><path d="M15 8C15.7956 8 16.5587 7.68393 17.1213 7.12132C17.6839 6.55871 18 5.79565 18 5C18 4.20435 17.6839 3.44129 17.1213 2.87868C16.5587 2.31607 15.7956 2 15 2" stroke-linecap="round" stroke-linejoin="round"/><path d="M4 21V19C4 17.9391 4.42143 16.9217 5.17157 16.1716C5.92172 15.4214 6.93913 15 8 15H10C11.0609 15 12.0783 15.4214 12.8284 16.1716C13.5786 16.9217 14 17.9391 14 19V21" stroke-linecap="round" stroke-linejoin="round"/><path d="M18 15C19.0609 15 20.0783 15.4214 20.8284 16.1716C21.5786 16.9217 22 17.9391 22 19V21" stroke-linecap="round" stroke-linejoin="round"/>`,

        // Edit/Pencil icon
        edit: `<path d="M11 4H4C3.46957 4 2.96086 4.21071 2.58579 4.58579C2.21071 4.96086 2 5.46957 2 6V20C2 20.5304 2.21071 21.0391 2.58579 21.4142C2.96086 21.7893 3.46957 22 4 22H18C18.5304 22 19.0391 21.7893 19.4142 21.4142C19.7893 21.0391 20 20.5304 20 20V13" stroke-linecap="round" stroke-linejoin="round"/><path d="M18.5 2.50001C18.8978 2.10219 19.4374 1.87868 20 1.87868C20.5626 1.87868 21.1022 2.10219 21.5 2.50001C21.8978 2.89784 22.1213 3.4374 22.1213 4.00001C22.1213 4.56262 21.8978 5.10219 21.5 5.50001L12 15L8 16L9 12L18.5 2.50001Z" stroke-linecap="round" stroke-linejoin="round"/>`,

        // Close/X icon
        close: `<path d="M18 6L6 18M6 6L18 18" stroke-linecap="round" stroke-linejoin="round"/>`,

        // Check/Success icon
        check: `<path d="M20 6L9 17L4 12" stroke-linecap="round" stroke-linejoin="round"/>`,

        // Warning/Alert icon
        warning: `<path d="M12 9V13M12 17H12.01M10.29 3.86L1.82 18C1.64537 18.3024 1.55297 18.6453 1.55199 18.9945C1.55101 19.3437 1.6415 19.6871 1.81443 19.9905C1.98737 20.2939 2.23673 20.5467 2.53835 20.724C2.83997 20.9012 3.18256 20.9968 3.53 21H20.47C20.8174 20.9968 21.16 20.9012 21.4617 20.724C21.7633 20.5467 22.0126 20.2939 22.1856 19.9905C22.3585 19.6871 22.449 19.3437 22.448 18.9945C22.447 18.6453 22.3546 18.3024 22.18 18L13.71 3.86C13.5318 3.56611 13.2807 3.32313 12.9812 3.15449C12.6817 2.98585 12.3438 2.89725 12 2.89725C11.6562 2.89725 11.3183 2.98585 11.0188 3.15449C10.7193 3.32313 10.4682 3.56611 10.29 3.86Z" stroke-linecap="round" stroke-linejoin="round"/>`,

        // Calendar icon
        calendar: `<rect x="3" y="4" width="18" height="18" rx="2" ry="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M16 2V6M8 2V6M3 10H21" stroke-linecap="round" stroke-linejoin="round"/>`,

        // Male icon
        male: `<circle cx="10" cy="8" r="5" stroke-linecap="round" stroke-linejoin="round"/><path d="M19 3L14 8M19 3H15M19 3V7" stroke-linecap="round" stroke-linejoin="round"/>`,

        // Female icon
        female: `<circle cx="12" cy="8" r="5" stroke-linecap="round" stroke-linejoin="round"/><path d="M12 13V21M9 18H15" stroke-linecap="round" stroke-linejoin="round"/>`,

        // Arrow left/back
        "arrow-left": `<path d="M19 12H5M12 19L5 12L12 5" stroke-linecap="round" stroke-linejoin="round"/>`,

        // Arrow right/next
        "arrow-right": `<path d="M5 12H19M12 5L19 12L12 19" stroke-linecap="round" stroke-linejoin="round"/>`,

        // Dashboard/Grid
        dashboard: `<rect x="3" y="3" width="7" height="7" rx="1" stroke-linecap="round" stroke-linejoin="round"/><rect x="14" y="3" width="7" height="7" rx="1" stroke-linecap="round" stroke-linejoin="round"/><rect x="14" y="14" width="7" height="7" rx="1" stroke-linecap="round" stroke-linejoin="round"/><rect x="3" y="14" width="7" height="7" rx="1" stroke-linecap="round" stroke-linejoin="round"/>`,

        // Settings/Cog
        settings: `<circle cx="12" cy="12" r="3" stroke-linecap="round" stroke-linejoin="round"/><path d="M19.4 15C19.2669 15.3016 19.2272 15.6362 19.286 15.9606C19.3448 16.285 19.4995 16.5843 19.73 16.82L19.79 16.88C19.976 17.0657 20.1235 17.2863 20.2241 17.5291C20.3248 17.7719 20.3766 18.0322 20.3766 18.295C20.3766 18.5578 20.3248 18.8181 20.2241 19.0609C20.1235 19.3037 19.976 19.5243 19.79 19.71C19.6043 19.896 19.3837 20.0435 19.1409 20.1441C18.8981 20.2448 18.6378 20.2966 18.375 20.2966C18.1122 20.2966 17.8519 20.2448 17.6091 20.1441C17.3663 20.0435 17.1457 19.896 16.96 19.71L16.9 19.65C16.6643 19.4195 16.365 19.2648 16.0406 19.206C15.7162 19.1472 15.3816 19.1869 15.08 19.32C14.7842 19.4468 14.532 19.6572 14.3543 19.9255C14.1766 20.1938 14.0813 20.5082 14.08 20.83V21C14.08 21.5304 13.8693 22.0391 13.4942 22.4142C13.1191 22.7893 12.6104 23 12.08 23C11.5496 23 11.0409 22.7893 10.6658 22.4142C10.2907 22.0391 10.08 21.5304 10.08 21V20.91C10.0723 20.579 9.96512 20.258 9.77251 19.9887C9.5799 19.7194 9.31074 19.5143 9 19.4C8.69838 19.2669 8.36381 19.2272 8.03941 19.286C7.71502 19.3448 7.41568 19.4995 7.18 19.73L7.12 19.79C6.93425 19.976 6.71368 20.1235 6.47088 20.2241C6.22808 20.3248 5.96783 20.3766 5.705 20.3766C5.44217 20.3766 5.18192 20.3248 4.93912 20.2241C4.69632 20.1235 4.47575 19.976 4.29 19.79C4.10405 19.6043 3.95653 19.3837 3.85588 19.1409C3.75523 18.8981 3.70343 18.6378 3.70343 18.375C3.70343 18.1122 3.75523 17.8519 3.85588 17.6091C3.95653 17.3663 4.10405 17.1457 4.29 16.96L4.35 16.9C4.58054 16.6643 4.73519 16.365 4.794 16.0406C4.85282 15.7162 4.81312 15.3816 4.68 15.08C4.55324 14.7842 4.34276 14.532 4.07447 14.3543C3.80618 14.1766 3.49179 14.0813 3.17 14.08H3C2.46957 14.08 1.96086 13.8693 1.58579 13.4942C1.21071 13.1191 1 12.6104 1 12.08C1 11.5496 1.21071 11.0409 1.58579 10.6658C1.96086 10.2907 2.46957 10.08 3 10.08H3.09C3.42099 10.0723 3.74197 9.96512 4.0113 9.77251C4.28063 9.5799 4.48572 9.31074 4.6 9C4.73312 8.69838 4.77282 8.36381 4.714 8.03941C4.65519 7.71502 4.50054 7.41568 4.27 7.18L4.21 7.12C4.02405 6.93425 3.87653 6.71368 3.77588 6.47088C3.67523 6.22808 3.62343 5.96783 3.62343 5.705C3.62343 5.44217 3.67523 5.18192 3.77588 4.93912C3.87653 4.69632 4.02405 4.47575 4.21 4.29C4.39575 4.10405 4.61632 3.95653 4.85912 3.85588C5.10192 3.75523 5.36217 3.70343 5.625 3.70343C5.88783 3.70343 6.14808 3.75523 6.39088 3.85588C6.63368 3.95653 6.85425 4.10405 7.04 4.29L7.1 4.35C7.33568 4.58054 7.63502 4.73519 7.95941 4.794C8.28381 4.85282 8.61838 4.81312 8.92 4.68H9C9.29577 4.55324 9.54802 4.34276 9.72569 4.07447C9.90337 3.80618 9.99872 3.49179 10 3.17V3C10 2.46957 10.2107 1.96086 10.5858 1.58579C10.9609 1.21071 11.4696 1 12 1C12.5304 1 13.0391 1.21071 13.4142 1.58579C13.7893 1.96086 14 2.46957 14 3V3.09C14.0013 3.41179 14.0966 3.72618 14.2743 3.99447C14.452 4.26276 14.7042 4.47324 15 4.6C15.3016 4.73312 15.6362 4.77282 15.9606 4.714C16.285 4.65519 16.5843 4.50054 16.82 4.27L16.88 4.21C17.0657 4.02405 17.2863 3.87653 17.5291 3.77588C17.7719 3.67523 18.0322 3.62343 18.295 3.62343C18.5578 3.62343 18.8181 3.67523 19.0609 3.77588C19.3037 3.87653 19.5243 4.02405 19.71 4.21C19.896 4.39575 20.0435 4.61632 20.1441 4.85912C20.2448 5.10192 20.2966 5.36217 20.2966 5.625C20.2966 5.88783 20.2448 6.14808 20.1441 6.39088C20.0435 6.63368 19.896 6.85425 19.71 7.04L19.65 7.1C19.4195 7.33568 19.2648 7.63502 19.206 7.95941C19.1472 8.28381 19.1869 8.61838 19.32 8.92V9C19.4468 9.29577 19.6572 9.54802 19.9255 9.72569C20.1938 9.90337 20.5082 9.99872 20.83 10H21C21.5304 10 22.0391 10.2107 22.4142 10.5858C22.7893 10.9609 23 11.4696 23 12C23 12.5304 22.7893 13.0391 22.4142 13.4142C22.0391 13.7893 21.5304 14 21 14H20.91C20.5882 14.0013 20.2738 14.0966 20.0055 14.2743C19.7372 14.452 19.5268 14.7042 19.4 15Z" stroke-linecap="round" stroke-linejoin="round"/>`,

        // Home
        home: `<path d="M3 12L5 10M5 10L12 3L19 10M5 10V20C5 20.2652 5.10536 20.5196 5.29289 20.7071C5.48043 20.8946 5.73478 21 6 21H9M19 10L21 12M19 10V20C19 20.2652 18.8946 20.5196 18.7071 20.7071C18.5196 20.8946 18.2652 21 18 21H15M9 21C9.26522 21 9.51957 20.8946 9.70711 20.7071C9.89464 20.5196 10 20.2652 10 20V16C10 15.7348 10.1054 15.4804 10.2929 15.2929C10.4804 15.1054 10.7348 15 11 15H13C13.2652 15 13.5196 15.1054 13.7071 15.2929C13.8946 15.4804 14 15.7348 14 16V20C14 20.2652 14.1054 20.5196 14.2929 20.7071C14.4804 20.8946 14.7348 21 15 21M9 21H15" stroke-linecap="round" stroke-linejoin="round"/>`,

        // Refresh
        refresh: `<path d="M1 4V10H7" stroke-linecap="round" stroke-linejoin="round"/><path d="M23 20V14H17" stroke-linecap="round" stroke-linejoin="round"/><path d="M20.49 9C19.9828 7.56678 19.1209 6.28535 17.9845 5.27542C16.8482 4.26549 15.4745 3.56076 13.9917 3.22426C12.5089 2.88776 10.9652 2.93026 9.50481 3.34783C8.04437 3.76539 6.71475 4.54467 5.64 5.61L1 10M23 14L18.36 18.39C17.2853 19.4553 15.9556 20.2346 14.4952 20.6522C13.0348 21.0697 11.4911 21.1122 10.0083 20.7757C8.52547 20.4392 7.1518 19.7345 6.01547 18.7246C4.87913 17.7147 4.01717 16.4332 3.51 15" stroke-linecap="round" stroke-linejoin="round"/>`,

        // Download
        download: `<path d="M21 15V19C21 19.5304 20.7893 20.0391 20.4142 20.4142C20.0391 20.7893 19.5304 21 19 21H5C4.46957 21 3.96086 20.7893 3.58579 20.4142C3.21071 20.0391 3 19.5304 3 19V15" stroke-linecap="round" stroke-linejoin="round"/><path d="M7 10L12 15L17 10" stroke-linecap="round" stroke-linejoin="round"/><path d="M12 15V3" stroke-linecap="round" stroke-linejoin="round"/>`,

        // Filter
        filter: `<path d="M22 3H2L10 12.46V19L14 21V12.46L22 3Z" stroke-linecap="round" stroke-linejoin="round"/>`,

        // Menu (hamburger)
        menu: `<path d="M3 12H21M3 6H21M3 18H21" stroke-linecap="round" stroke-linejoin="round"/>`,

        // Logout
        logout: `<path d="M9 21H5C4.46957 21 3.96086 20.7893 3.58579 20.4142C3.21071 20.0391 3 19.5304 3 19V5C3 4.46957 3.21071 3.96086 3.58579 3.58579C3.96086 3.21071 4.46957 3 5 3H9" stroke-linecap="round" stroke-linejoin="round"/><path d="M16 17L21 12L16 7" stroke-linecap="round" stroke-linejoin="round"/><path d="M21 12H9" stroke-linecap="round" stroke-linejoin="round"/>`,

        // Lock
        lock: `<rect x="3" y="11" width="18" height="11" rx="2" ry="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M7 11V7C7 5.67392 7.52678 4.40215 8.46447 3.46447C9.40215 2.52678 10.6739 2 12 2C13.3261 2 14.5979 2.52678 15.5355 3.46447C16.4732 4.40215 17 5.67392 17 7V11" stroke-linecap="round" stroke-linejoin="round"/>`,

        // Book/Subject icon
        book: `<path d="M4 19.5C4 18.837 4.26339 18.2011 4.73223 17.7322C5.20107 17.2634 5.83696 17 6.5 17H20" stroke-linecap="round" stroke-linejoin="round"/><path d="M6.5 2H20V22H6.5C5.83696 22 5.20107 21.7366 4.73223 21.2678C4.26339 20.7989 4 20.163 4 19.5V4.5C4 3.83696 4.26339 3.20107 4.73223 2.73223C5.20107 2.26339 5.83696 2 6.5 2Z" stroke-linecap="round" stroke-linejoin="round"/>`,

        // Graduation/Student icon
        graduation: `<path d="M22 10v6M2 10l10-5 10 5-10 5z" stroke-linecap="round" stroke-linejoin="round"/><path d="M6 12v5c3 3 9 3 12 0v-5" stroke-linecap="round" stroke-linejoin="round"/>`,

        // Activity/Flag icon
        activity: `<path d="M4 15s1-1 4-1 5 2 8 2 4-1 4-1V3s-1 1-4 1-5-2-8-2-4 1-4 1z" stroke-linecap="round" stroke-linejoin="round"/><path d="M4 22v-7" stroke-linecap="round" stroke-linejoin="round"/>`,
    };
</script>

<svg
    width={size}
    height={size}
    viewBox="0 0 24 24"
    {fill}
    stroke={effectiveStroke}
    stroke-width={strokeWidth}
    class="icon icon-{name}"
>
    {@html icons[name] || ""}
</svg>

<style>
    svg {
        display: inline-block;
        vertical-align: middle;
        flex-shrink: 0;
    }
</style>
