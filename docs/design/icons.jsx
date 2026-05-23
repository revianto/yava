// YAVA — Tabler-style outline icons. All inherit currentColor.
// Sizes via the .ico / .ico--lg / .ico--xl classes or style overrides.

const Icon = ({ children, size = 18, className = "", strokeWidth = 1.75, style }) => (
  <svg
    className={"ico " + className}
    width={size}
    height={size}
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth={strokeWidth}
    strokeLinecap="round"
    strokeLinejoin="round"
    style={style}
  >
    {children}
  </svg>
);

const IconPlay   = (p) => <Icon {...p}><path d="M7 4v16l13 -8z" /></Icon>;
const IconPause  = (p) => <Icon {...p}><rect x="6" y="5" width="4" height="14" rx="1" /><rect x="14" y="5" width="4" height="14" rx="1" /></Icon>;
const IconSkip   = (p) => <Icon {...p}><path d="M4 5v14l12 -7z" /><path d="M20 5v14" /></Icon>;
const IconReset  = (p) => <Icon {...p}><path d="M20 11a8 8 0 1 0 -2.34 5.66" /><path d="M20 4v7h-7" /></Icon>;
const IconTimer  = (p) => <Icon {...p}><circle cx="12" cy="13" r="8" /><path d="M9 3h6" /><path d="M12 8v5l3 2" /></Icon>;
const IconPlus   = (p) => <Icon {...p}><path d="M12 5v14M5 12h14" /></Icon>;
const IconSearch = (p) => <Icon {...p}><circle cx="11" cy="11" r="7" /><path d="m21 21 -4.3 -4.3" /></Icon>;
const IconBell   = (p) => <Icon {...p}><path d="M10 5a2 2 0 1 1 4 0 7 7 0 0 1 4 6v3a4 4 0 0 0 2 3H4a4 4 0 0 0 2 -3v-3a7 7 0 0 1 4 -6" /><path d="M9 17v1a3 3 0 0 0 6 0v-1" /></Icon>;
const IconStar   = (p) => <Icon {...p}><path d="m12 17.27 -5.18 3.36 1.64 -6.03 -4.46 -3.86 6.13 -.5L12 5l1.87 5.24 6.13 .5 -4.46 3.86 1.64 6.03z" /></Icon>;
const IconGroup  = (p) => <Icon {...p}><circle cx="9" cy="7" r="3" /><circle cx="17" cy="9" r="2" /><path d="M3 21v-2a4 4 0 0 1 4 -4h4a4 4 0 0 1 4 4v2" /><path d="M17 13a3 3 0 0 1 3 3v1" /></Icon>;
const IconArrow  = (p) => <Icon {...p}><path d="M5 12h14M13 6l6 6 -6 6" /></Icon>;
const IconArrowLeft = (p) => <Icon {...p}><path d="M19 12H5M11 6l-6 6 6 6" /></Icon>;
const IconArrowDown = (p) => <Icon {...p}><path d="M12 5v14M6 13l6 6 6 -6" /></Icon>;
const IconCheck  = (p) => <Icon {...p}><path d="M5 12l5 5L20 7" /></Icon>;
const IconClose  = (p) => <Icon {...p}><path d="M6 6l12 12M18 6 6 18" /></Icon>;
const IconEdit   = (p) => <Icon {...p}><path d="M4 20h4l10 -10 -4 -4L4 16v4z" /><path d="m13.5 6.5 4 4" /></Icon>;
const IconShare  = (p) => <Icon {...p}><circle cx="6" cy="12" r="3" /><circle cx="18" cy="6" r="3" /><circle cx="18" cy="18" r="3" /><path d="m8.6 13.5 6.8 3.5M8.6 10.5l6.8 -3.5" /></Icon>;
const IconHeat   = (p) => <Icon {...p}><path d="M12 3a4 4 0 0 0 -4 4c0 1.5 1 2.5 1.5 3.5C10 11.5 10 13 8 13a4 4 0 0 0 -4 4c0 2.5 2 4 4 4h8a4 4 0 0 0 4 -4c0 -2 -1 -3 -2 -4 -1.5 -1.5 -2 -3 -1 -4.5C18 7 18 3 12 3z" /></Icon>;
const IconDose   = (p) => <Icon {...p}><path d="M4 4h12l-2 16H6L4 4z" /><path d="M16 7h2a3 3 0 0 1 0 6h-2.5" /></Icon>;
const IconRatio  = (p) => <Icon {...p}><path d="M4 20 20 4" /><circle cx="6.5" cy="6.5" r="2.5" /><circle cx="17.5" cy="17.5" r="2.5" /></Icon>;
const IconYield  = (p) => <Icon {...p}><path d="M6 3h12l-1 6a5 5 0 0 1 -10 0L6 3z" /><path d="M12 14v4M9 21h6" /></Icon>;
const IconGrind  = (p) => <Icon {...p}><circle cx="12" cy="12" r="8" /><circle cx="12" cy="12" r="3" /><path d="M12 4v3M12 17v3M4 12h3M17 12h3" /></Icon>;
const IconKebab  = (p) => <Icon {...p}><circle cx="12" cy="5" r="1" /><circle cx="12" cy="12" r="1" /><circle cx="12" cy="19" r="1" /></Icon>;
const IconBookmark = (p) => <Icon {...p}><path d="M6 4h12v17l-6 -4 -6 4V4z" /></Icon>;

Object.assign(window, {
  Icon,
  IconPlay, IconPause, IconSkip, IconReset, IconTimer, IconPlus,
  IconSearch, IconBell, IconStar, IconGroup, IconArrow, IconArrowLeft, IconArrowDown,
  IconCheck, IconClose, IconEdit, IconShare,
  IconHeat, IconDose, IconRatio, IconYield, IconGrind, IconKebab, IconBookmark,
});
