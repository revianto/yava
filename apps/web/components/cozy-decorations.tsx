// Cozy decorative SVG elements — faceless figures, plants, coffee motifs.
// Each component is absolutely positioned; the parent must be position:relative (overflow:hidden recommended).

const TAN = '#E1BF91'
const BROWN = '#AD8257'
const CREAM = '#FBF8F1'

export function CozyFigureMug() {
  return (
    <div
      className="cozy-deco"
      aria-hidden="true"
      style={{ top: 24, right: 36, width: 200, height: 240, opacity: 0.42 }}
    >
      <svg viewBox="0 0 200 240" xmlns="http://www.w3.org/2000/svg" fill="none"
        stroke={TAN} strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
        <ellipse cx="100" cy="48" rx="26" ry="30"/>
        <path d="M76 38 Q88 18 100 18 Q114 18 124 36" strokeWidth="2"/>
        <path d="M55 130 Q56 96 90 86 L110 86 Q144 96 145 130 L150 230"/>
        <path d="M50 230 L55 130"/>
        <path d="M76 138 L76 188 Q76 196 84 196 L116 196 Q124 196 124 188 L124 138 Z"/>
        <path d="M124 150 Q140 152 140 168 Q140 184 124 184"/>
        <path d="M73 132 Q60 158 76 188"/>
        <path d="M127 132 Q140 158 124 188"/>
        <path d="M86 130 Q92 118 86 104 Q80 92 86 80" strokeWidth="2" opacity="0.7"/>
        <path d="M100 130 Q106 118 100 104" strokeWidth="2" opacity="0.7"/>
        <path d="M114 130 Q108 118 114 104 Q120 92 114 80" strokeWidth="2" opacity="0.7"/>
        <path d="M82 96 Q100 104 118 96" strokeWidth="2"/>
        <path d="M100 96 L100 130" strokeWidth="1.5" strokeDasharray="2 4" opacity="0.6"/>
      </svg>
    </div>
  )
}

export function CozyMugSteam() {
  return (
    <div
      className="cozy-deco"
      aria-hidden="true"
      style={{ top: 20, right: 20, width: 84, height: 100, opacity: 0.35 }}
    >
      <svg viewBox="0 0 120 140" xmlns="http://www.w3.org/2000/svg" fill="none"
        stroke={TAN} strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
        <path d="M22 60 L22 110 Q22 122 34 122 L70 122 Q82 122 82 110 L82 60 Z"/>
        <path d="M22 60 L82 60" strokeWidth="2"/>
        <path d="M82 72 Q102 75 102 88 Q102 102 82 105"/>
        <ellipse cx="52" cy="126" rx="42" ry="5"/>
        <path d="M34 50 Q40 38 34 24 Q28 12 34 0" opacity="0.7"/>
        <path d="M52 50 Q58 38 52 24" opacity="0.7"/>
        <path d="M70 50 Q64 38 70 24 Q76 12 70 0" opacity="0.7"/>
      </svg>
    </div>
  )
}

export function CozyBranch() {
  return (
    <div
      className="cozy-deco"
      aria-hidden="true"
      style={{ bottom: 20, right: 16, width: 140, height: 140, opacity: 0.85 }}
    >
      <svg viewBox="0 0 160 160" xmlns="http://www.w3.org/2000/svg" fill="none"
        stroke={TAN} strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
        <path d="M10 150 Q40 120 70 90 Q100 60 140 30"/>
        <path d="M55 110 Q42 100 38 86 Q52 92 60 104 Z" fill={TAN} fillOpacity="0.18"/>
        <path d="M48 88 L58 106" strokeWidth="1.5" opacity="0.6"/>
        <path d="M90 78 Q76 70 70 56 Q86 60 96 72 Z" fill={TAN} fillOpacity="0.18"/>
        <path d="M80 60 L92 76" strokeWidth="1.5" opacity="0.6"/>
        <circle cx="64" cy="96" r="8" fill={TAN} fillOpacity="0.75"/>
        <circle cx="78" cy="86" r="7" fill={TAN} fillOpacity="0.55"/>
        <circle cx="112" cy="52" r="9" fill={TAN} fillOpacity="0.85"/>
        <circle cx="126" cy="42" r="7" fill={TAN} fillOpacity="0.6"/>
        <path d="M64 88 L62 80" strokeWidth="1.5"/>
        <path d="M112 43 L108 36" strokeWidth="1.5"/>
      </svg>
    </div>
  )
}

export function CozyPlants() {
  return (
    <div
      className="cozy-deco"
      aria-hidden="true"
      style={{ top: -40, right: 200, width: 130, height: 170, opacity: 0.55 }}
    >
      <svg viewBox="0 0 160 200" xmlns="http://www.w3.org/2000/svg" fill="none"
        stroke={BROWN} strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <path d="M80 200 L80 130" strokeWidth="2.5"/>
        <path d="M80 140 Q72 95 64 60 Q72 38 80 25 Q88 38 96 60 Q88 95 80 140 Z"
          fill={BROWN} fillOpacity="0.12"/>
        <path d="M80 130 L80 30" strokeWidth="1.5" opacity="0.7"/>
        <path d="M80 150 Q56 120 36 92 Q34 70 46 60 Q60 80 76 122 Z"
          fill={BROWN} fillOpacity="0.10" transform="translate(-2 0)"/>
        <path d="M78 148 Q60 120 44 80" strokeWidth="1.5" opacity="0.6"/>
        <path d="M80 150 Q104 120 124 92 Q126 70 114 60 Q100 80 84 122 Z"
          fill={BROWN} fillOpacity="0.10" transform="translate(2 0)"/>
        <path d="M82 148 Q100 120 116 80" strokeWidth="1.5" opacity="0.6"/>
        <path d="M64 200 L96 200" strokeWidth="2.5"/>
        <path d="M68 200 L72 210" strokeWidth="2" opacity="0.5"/>
        <path d="M92 200 L88 210" strokeWidth="2" opacity="0.5"/>
      </svg>
    </div>
  )
}

export function CozyFigureWalking() {
  return (
    <div
      className="cozy-deco"
      aria-hidden="true"
      style={{ bottom: 16, right: 16, width: 160, height: 230, opacity: 0.35 }}
    >
      <svg viewBox="0 0 200 280" xmlns="http://www.w3.org/2000/svg" fill="none"
        stroke={CREAM} strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
        <ellipse cx="110" cy="40" rx="22" ry="26"/>
        <circle cx="92" cy="38" r="10" fill={CREAM} fillOpacity="0.25"/>
        <path d="M104 66 L104 76"/>
        <path d="M118 66 L118 76"/>
        <path d="M82 100 Q88 84 110 80 Q132 84 138 100 L142 180 Q124 192 100 192 Q80 192 76 180 Z"/>
        <path d="M138 110 Q156 130 148 158"/>
        <path d="M138 158 L160 158 L156 192 L142 192 Z"/>
        <path d="M140 168 L158 168" strokeWidth="1.5"/>
        <path d="M82 110 Q66 130 72 156"/>
        <path d="M96 192 L88 268"/>
        <path d="M124 192 L130 268"/>
        <path d="M84 268 L96 268"/>
        <path d="M126 268 L138 268"/>
        <path d="M148 152 Q152 142 148 132" strokeWidth="2" opacity="0.7"/>
        <path d="M156 152 Q160 142 156 132" strokeWidth="2" opacity="0.7"/>
      </svg>
    </div>
  )
}
