// YAVA — shared UI atoms.

const Logo = ({ size = 26, dark = false }) => (
  <a className="logo" href="#" onClick={(e) => e.preventDefault()}
    style={{ fontSize: size, color: dark ? "#fff" : "var(--deep-ink)" }}>
    <span>YAVA</span><span className="dot">.</span>
  </a>
);

// Tag — matches design system pill tags
const Tag = ({ children, variant = "default", size = "sm", style }) => (
  <span className={`tag tag--${variant} ${size === "lg" ? "tag--lg" : ""}`} style={style}>{children}</span>
);

// Map recipe tag string → variant
const tagVariant = (label) => {
  const k = String(label).toUpperCase();
  if (k === "ESPRESSO") return "espresso";
  if (k === "V60") return "v60";
  if (k.includes("COLD")) return "cold";
  if (k === "GRUP") return "group";
  return "default";
};

// Visibility tag
const VisibilityTag = ({ visibility, isDefault }) => {
  if (isDefault) return <Tag variant="default">DEFAULT</Tag>;
  if (visibility === "public") return <Tag variant="public">PUBLIK</Tag>;
  if (visibility === "group") return <Tag variant="group">GRUP</Tag>;
  return <Tag variant="private">PRIBADI</Tag>;
};

// Avatar
const Avatar = ({ initials, color, size = "md" }) => (
  <span className={`avatar ${size === "sm" ? "avatar--sm" : ""}`} style={color ? { background: color } : null}>
    {initials}
  </span>
);

// Timer format mm:ss
const fmt = (seconds) => {
  const s = Math.max(0, Math.round(seconds));
  const m = Math.floor(s / 60);
  const r = s % 60;
  return `${String(m).padStart(2, "0")}:${String(r).padStart(2, "0")}`;
};

// Total duration of a recipe
const totalDuration = (recipe) =>
  (recipe.timeline || []).filter((s) => s.kind === "session").reduce((a, s) => a + s.duration, 0);
const sessionCount = (recipe) => (recipe.timeline || []).filter((s) => s.kind === "session").length;

// Param row (5 cells)
const Params = ({ params }) => (
  <div className="params">
    <div className="params__cell"><span className="lbl">Dose</span><span className="val t-mono-num">{params.dose}</span></div>
    <div className="params__cell"><span className="lbl">Yield</span><span className="val t-mono-num">{params.yield}</span></div>
    <div className="params__cell"><span className="lbl">Suhu</span><span className="val t-mono-num">{params.temp}</span></div>
    <div className="params__cell"><span className="lbl">Grind</span><span className="val">{params.grind}</span></div>
    <div className="params__cell"><span className="lbl">Ratio</span><span className="val t-mono-num">{params.ratio}</span></div>
  </div>
);

// Recipe Card (white, used in dashboard grid)
const RecipeCard = ({ recipe, onClick, accent = false }) => {
  const total = totalDuration(recipe);
  const sessions = sessionCount(recipe);
  return (
    <button className="recipe-card" onClick={onClick}>
      <div className="row gap-2 wrap" style={{ marginBottom: 2 }}>
        {recipe.tags.slice(0, 2).map((t) => (
          <Tag key={t} variant={tagVariant(t)}>{t}</Tag>
        ))}
        <span style={{ flex: 1 }} />
        <VisibilityTag visibility={recipe.visibility} isDefault={recipe.isDefault} />
      </div>
      <div className="recipe-card__title">{recipe.name}</div>
      <div className="recipe-card__meta">
        <span>{sessions} sesi</span>
        <span className="dot" />
        <span className="t-mono-num">{Math.floor(total / 60)}m {total % 60}s</span>
        <span className="dot" />
        <span>{recipe.params.ratio}</span>
      </div>
      <div className="row between" style={{ marginTop: "auto", paddingTop: 12 }}>
        <span className="t-caption">{recipe.subtype}</span>
        <span className="row gap-1" style={{ color: "var(--deep-ink)", fontWeight: 700, fontSize: 13 }}>
          Lihat <IconArrow size={14} />
        </span>
      </div>
    </button>
  );
};

// Hero recipe card (dark surface) for top of dashboard
const HeroRecipeCard = ({ recipe, onStart, onOpen }) => {
  const total = totalDuration(recipe);
  const sessions = sessionCount(recipe);
  return (
    <div className="card card--dark card--hero" style={{ position: "relative", overflow: "hidden", padding: 40 }}>
      {/* Soft background mark */}
      <div aria-hidden style={{
        position: "absolute", right: -60, top: -60, width: 320, height: 320, borderRadius: "50%",
        background: "radial-gradient(circle, rgba(61,43,255,.45), transparent 70%)", filter: "blur(0px)"
      }} />
      <div className="row gap-2 wrap" style={{ marginBottom: 18, position: "relative" }}>
        <Tag variant="espresso" size="lg">RESEP MINGGU INI</Tag>
        {recipe.tags.map((t) => (
          <Tag key={t} variant={tagVariant(t) === "default" ? "private" : tagVariant(t)} size="lg">{t}</Tag>
        ))}
      </div>
      <div style={{ position: "relative", display: "grid", gridTemplateColumns: "1.4fr 1fr", gap: 48, alignItems: "end" }}>
        <div>
          <div className="t-display" style={{ marginBottom: 16, maxWidth: 540 }}>{recipe.name}</div>
          <div className="t-body muted-dark" style={{ maxWidth: 520, marginBottom: 28 }}>{recipe.description}</div>
          <div className="row gap-2">
            <button className="btn btn--primary btn--lg" onClick={onStart}>
              <IconPlay size={16} /> Mulai Brewing
            </button>
            <button className="btn btn--secondary-dark btn--lg" onClick={onOpen}>Lihat resep</button>
          </div>
        </div>
        <div className="col gap-3" style={{ alignItems: "stretch" }}>
          <div className="row between" style={{ paddingBottom: 14, borderBottom: "1px solid rgba(255,255,255,.10)" }}>
            <span className="t-label muted-dark">Total brew</span>
            <span className="t-h2 t-mono-num">{Math.floor(total / 60)}:{String(total % 60).padStart(2, "0")}</span>
          </div>
          <div className="row between" style={{ paddingBottom: 14, borderBottom: "1px solid rgba(255,255,255,.10)" }}>
            <span className="t-label muted-dark">Sesi</span>
            <span className="t-h2 t-mono-num">{sessions}</span>
          </div>
          <div className="row between" style={{ paddingBottom: 14, borderBottom: "1px solid rgba(255,255,255,.10)" }}>
            <span className="t-label muted-dark">Ratio</span>
            <span className="t-h2">{recipe.params.ratio}</span>
          </div>
          <div className="row between">
            <span className="t-label muted-dark">Suhu</span>
            <span className="t-h2 t-mono-num">{recipe.params.temp}</span>
          </div>
        </div>
      </div>
    </div>
  );
};

// Topnav shared across screens
const Topnav = ({ active, onNav, onBrew, onNew }) => (
  <header className="topnav">
    <div className="row gap-4">
      <Logo size={26} />
      <nav className="topnav__links">
        {["Dashboard", "Resep saya", "Explore", "Grup"].map((label, i) => {
          const key = ["dashboard", "library", "explore", "groups"][i];
          return (
            <a key={key}
               href="#"
               onClick={(e) => { e.preventDefault(); onNav && onNav(key); }}
               className={`topnav__link ${active === key ? "topnav__link--active" : ""}`}>
              {label}
            </a>
          );
        })}
      </nav>
    </div>
    <div className="row gap-2">
      <div style={{ position: "relative" }}>
        <IconSearch size={18} style={{ position: "absolute", left: 14, top: "50%", transform: "translateY(-50%)", color: "var(--muted)" }} />
        <input className="input input--search" placeholder="Cari resep, jenis, atau bahan…" style={{ width: 280 }} />
      </div>
      <button className="icon-btn" title="Notifikasi"><IconBell size={18} /></button>
      <button className="btn btn--light-primary" onClick={onNew}><IconPlus size={16} /> Buat resep</button>
      <Avatar initials="ND" color="var(--coral-red)" />
    </div>
  </header>
);

// Section header
const SectionHeader = ({ kicker, title, action }) => (
  <div className="row between" style={{ marginBottom: 16, alignItems: "flex-end" }}>
    <div>
      {kicker && <div className="t-label" style={{ color: "var(--muted)", marginBottom: 6 }}>{kicker}</div>}
      <div className="t-h1">{title}</div>
    </div>
    {action}
  </div>
);

// Timeline step row (used on Recipe Detail)
const StepRow = ({ index, step }) => {
  if (step.kind === "note") {
    return (
      <div className="step">
        <div className="step__num step__num--note">N</div>
        <div>
          <div className="step__name muted">Catatan</div>
          <div className="step__note">{step.content}</div>
        </div>
        <div className="step__time muted">—</div>
      </div>
    );
  }
  return (
    <div className="step">
      <div className="step__num">{index}</div>
      <div>
        <div className="step__name">{step.name}</div>
        <div className="step__note">{step.note}</div>
      </div>
      <div className="step__time">{step.duration}<small>s</small></div>
    </div>
  );
};

Object.assign(window, {
  Logo, Tag, tagVariant, VisibilityTag, Avatar,
  Params, RecipeCard, HeroRecipeCard, Topnav, SectionHeader, StepRow,
  fmt, totalDuration, sessionCount,
});
