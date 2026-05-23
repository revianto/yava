import { SVGProps } from 'react'

interface IconProps extends SVGProps<SVGSVGElement> {
  size?: number
  strokeWidth?: number
}

const Icon = ({ children, size = 18, strokeWidth = 1.75, className = '', ...rest }: IconProps & { children: React.ReactNode }) => (
  <svg
    className={'ico ' + className}
    width={size}
    height={size}
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth={strokeWidth}
    strokeLinecap="round"
    strokeLinejoin="round"
    {...rest}
  >
    {children}
  </svg>
)

export const IconPlay     = (p: IconProps) => <Icon {...p}><path d="M7 4v16l13 -8z" /></Icon>
export const IconPause    = (p: IconProps) => <Icon {...p}><rect x="6" y="5" width="4" height="14" rx="1" /><rect x="14" y="5" width="4" height="14" rx="1" /></Icon>
export const IconSkip     = (p: IconProps) => <Icon {...p}><path d="M4 5v14l12 -7z" /><path d="M20 5v14" /></Icon>
export const IconReset    = (p: IconProps) => <Icon {...p}><path d="M20 11a8 8 0 1 0 -2.34 5.66" /><path d="M20 4v7h-7" /></Icon>
export const IconPlus     = (p: IconProps) => <Icon {...p}><path d="M12 5v14M5 12h14" /></Icon>
export const IconSearch   = (p: IconProps) => <Icon {...p}><circle cx="11" cy="11" r="7" /><path d="m21 21 -4.3 -4.3" /></Icon>
export const IconBell     = (p: IconProps) => <Icon {...p}><path d="M10 5a2 2 0 1 1 4 0 7 7 0 0 1 4 6v3a4 4 0 0 0 2 3H4a4 4 0 0 0 2 -3v-3a7 7 0 0 1 4 -6" /><path d="M9 17v1a3 3 0 0 0 6 0v-1" /></Icon>
export const IconStar     = (p: IconProps) => <Icon {...p}><path d="m12 17.27 -5.18 3.36 1.64 -6.03 -4.46 -3.86 6.13 -.5L12 5l1.87 5.24 6.13 .5 -4.46 3.86 1.64 6.03z" /></Icon>
export const IconGroup    = (p: IconProps) => <Icon {...p}><circle cx="9" cy="7" r="3" /><circle cx="17" cy="9" r="2" /><path d="M3 21v-2a4 4 0 0 1 4 -4h4a4 4 0 0 1 4 4v2" /><path d="M17 13a3 3 0 0 1 3 3v1" /></Icon>
export const IconArrow    = (p: IconProps) => <Icon {...p}><path d="M5 12h14M13 6l6 6 -6 6" /></Icon>
export const IconArrowLeft  = (p: IconProps) => <Icon {...p}><path d="M19 12H5M11 6l-6 6 6 6" /></Icon>
export const IconArrowDown  = (p: IconProps) => <Icon {...p}><path d="M12 5v14M6 13l6 6 6 -6" /></Icon>
export const IconCheck    = (p: IconProps) => <Icon {...p}><path d="M5 12l5 5L20 7" /></Icon>
export const IconClose    = (p: IconProps) => <Icon {...p}><path d="M6 6l12 12M18 6 6 18" /></Icon>
export const IconEdit     = (p: IconProps) => <Icon {...p}><path d="M4 20h4l10 -10 -4 -4L4 16v4z" /><path d="m13.5 6.5 4 4" /></Icon>
export const IconShare    = (p: IconProps) => <Icon {...p}><circle cx="6" cy="12" r="3" /><circle cx="18" cy="6" r="3" /><circle cx="18" cy="18" r="3" /><path d="m8.6 13.5 6.8 3.5M8.6 10.5l6.8 -3.5" /></Icon>
export const IconKebab    = (p: IconProps) => <Icon {...p}><circle cx="12" cy="5" r="1" /><circle cx="12" cy="12" r="1" /><circle cx="12" cy="19" r="1" /></Icon>
export const IconBookmark = (p: IconProps) => <Icon {...p}><path d="M6 4h12v17l-6 -4 -6 4V4z" /></Icon>
