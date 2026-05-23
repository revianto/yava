// YAVA — Dashboard screen
// Hero recipe (dark) + tabbed library + recent activity.

const Dashboard = ({ onOpenRecipe, onStart, onNav, onNew, tweak }) => {
  const { RECIPES, TYPES, HERO_RECIPE_ID } = window.YAVA_DATA;
  const hero = RECIPES.find((r) => r.id === HERO_RECIPE_ID);
  const [activeType, setActiveType] = React.useState("Semua");

  const filtered = RECIPES.filter((r) => activeType === "Semua" || r.type === activeType);

  return (
    <div className="page fade-in">
      <Topnav active="dashboard" onNav={onNav} onNew={onNew} />

      <main className="container col gap-4" style={{ paddingTop: 16, gap: 48 }}>
        {/* Greeting + hero */}
        <section className="col gap-3">
          <div className="row between" style={{ alignItems: "flex-end" }}>
            <div>
              <div className="t-label" style={{ color: "var(--muted)", marginBottom: 8 }}>Selasa · 23 Mei 2026</div>
              <div className="t-display" style={{ fontSize: 48, letterSpacing: "-.025em" }}>
                Pagi, Nadira.<span style={{ color: "var(--coral-red)" }}>.</span>
              </div>
              <div className="t-body muted" style={{ marginTop: 6 }}>
                3 resep aktif · 12 sesi brewing minggu ini.
              </div>
            </div>
            <div className="row gap-2">
              <button className="btn btn--light-secondary"><IconBookmark size={16} /> Favorit</button>
              <button className="btn btn--light-primary" onClick={onNew}><IconPlus size={16} /> Buat resep</button>
            </div>
          </div>

          <HeroRecipeCard recipe={hero} onStart={() => onStart(hero.id)} onOpen={() => onOpenRecipe(hero.id)} />
        </section>

        {/* Library */}
        <section className="col gap-3">
          <div className="row between" style={{ alignItems: "flex-end" }}>
            <div>
              <div className="t-label" style={{ color: "var(--muted)", marginBottom: 6 }}>Koleksi</div>
              <div className="t-h1">Resep saya</div>
            </div>
            <div className="tabs">
              {TYPES.map((t) => (
                <button
                  key={t}
                  className={`tab ${t === activeType ? "tab--active" : ""}`}
                  onClick={() => setActiveType(t)}
                >
                  {t}
                </button>
              ))}
            </div>
          </div>

          <div style={{ display: "grid", gridTemplateColumns: "repeat(3, 1fr)", gap: 16 }}>
            {filtered.map((r) => (
              <RecipeCard key={r.id} recipe={r} onClick={() => onOpenRecipe(r.id)} />
            ))}
            {/* New recipe slot */}
            <button
              className="recipe-card"
              onClick={onNew}
              style={{
                background: "transparent",
                border: "1.5px dashed var(--deep-ink)",
                alignItems: "center",
                justifyContent: "center",
                minHeight: 200,
                color: "var(--deep-ink)",
                fontWeight: 700,
              }}
            >
              <IconPlus size={28} />
              <span style={{ marginTop: 4 }}>Buat resep baru</span>
              <span className="t-caption" style={{ color: "var(--muted)" }}>
                Tambahkan sesi timer atau notes
              </span>
            </button>
          </div>
        </section>

        {/* Two-up: Activity + Group */}
        <section style={{ display: "grid", gridTemplateColumns: "1.4fr 1fr", gap: 16 }}>
          {/* Activity */}
          <div className="card" style={{ padding: 24 }}>
            <div className="row between" style={{ marginBottom: 16 }}>
              <div className="t-h2">Aktivitas terakhir</div>
              <a href="#" onClick={(e) => e.preventDefault()} className="t-small" style={{ color: "var(--deep-ink)", fontWeight: 700 }}>
                Lihat semua
              </a>
            </div>
            <div className="col">
              {[
                { time: "08:14", action: "Selesai brewing", what: "V60 Light Roast 15g / 250ml", icon: <IconCheck size={16} /> },
                { time: "Kemarin", action: "Edit resep", what: "Flair Espresso 18g / 36g", icon: <IconEdit size={16} /> },
                { time: "Senin", action: "Disetujui di grup", what: "Kopi Susu Aren ke Komunitas Senayan", icon: <IconGroup size={16} /> },
                { time: "Minggu", action: "Duplikat", what: "V60 Dark Roast 18g / 270ml dari Sistem YAVA", icon: <IconShare size={16} /> },
              ].map((row, i) => (
                <div key={i} className="row gap-3" style={{ padding: "14px 0", borderTop: i === 0 ? "0" : "1px solid var(--hairline)" }}>
                  <div style={{
                    width: 32, height: 32, borderRadius: "50%",
                    display: "inline-flex", alignItems: "center", justifyContent: "center",
                    background: "var(--lavender-fog)", color: "var(--deep-ink)"
                  }}>{row.icon}</div>
                  <div className="grow">
                    <div style={{ fontWeight: 700, fontSize: 14 }}>{row.action}</div>
                    <div className="t-small muted">{row.what}</div>
                  </div>
                  <div className="t-small muted t-mono-num">{row.time}</div>
                </div>
              ))}
            </div>
          </div>

          {/* Group highlight (Electric purple) */}
          <div className="card card--electric card--hero" style={{ padding: 28, position: "relative" }}>
            <Tag variant="public" size="lg" style={{ background: "rgba(255,255,255,.12)", color: "var(--lilac)", border: "0" }}>
              GRUP AKTIF
            </Tag>
            <div className="t-h1" style={{ color: "#fff", marginTop: 16, lineHeight: 1.05 }}>
              Komunitas<br />Senayan
            </div>
            <div className="t-small" style={{ color: "var(--lilac)", marginTop: 8, marginBottom: 24 }}>
              28 anggota · 42 resep aktif · 2 menunggu approval
            </div>
            <div className="row gap-1" style={{ marginBottom: 20 }}>
              {["RW", "AD", "MS", "+25"].map((n, i) => (
                <span key={i} className="avatar avatar--sm" style={{
                  background: i === 3 ? "rgba(255,255,255,.18)" : "var(--lilac)",
                  color: i === 3 ? "var(--lilac)" : "var(--electric)",
                  border: "2px solid var(--electric)",
                  marginLeft: i === 0 ? 0 : -8,
                }}>{n}</span>
              ))}
            </div>
            <button className="btn" style={{ background: "var(--lilac)", color: "var(--electric)" }}>
              Buka grup <IconArrow size={14} />
            </button>
          </div>
        </section>
      </main>
    </div>
  );
};

window.Dashboard = Dashboard;
