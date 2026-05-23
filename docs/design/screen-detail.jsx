// YAVA — Recipe Detail screen
// Recipe params + interleaved sessions/notes + big start CTA.

const RecipeDetail = ({ recipeId, onBack, onStart, onNav }) => {
  const recipe = window.YAVA_DATA.RECIPES.find((r) => r.id === recipeId);
  if (!recipe) return <div className="page" style={{ padding: 48 }}>Resep tidak ditemukan.</div>;

  const total = totalDuration(recipe);
  const sessions = sessionCount(recipe);

  let sessionIdx = 0;

  return (
    <div className="page fade-in">
      <Topnav active="library" onNav={onNav} />

      <main className="container">
        {/* Breadcrumb */}
        <button
          onClick={onBack}
          style={{
            background: "transparent", border: 0, cursor: "pointer",
            display: "inline-flex", alignItems: "center", gap: 6,
            color: "var(--muted)", fontWeight: 500, fontSize: 13, padding: 0, marginBottom: 24,
          }}
        >
          <IconArrowLeft size={14} /> Kembali ke {recipe.type} · {recipe.subtype}
        </button>

        {/* Title block */}
        <div className="row between" style={{ alignItems: "flex-start", gap: 32, marginBottom: 32 }}>
          <div style={{ maxWidth: 720 }}>
            <div className="row gap-2 wrap" style={{ marginBottom: 14 }}>
              {recipe.tags.map((t) => (
                <Tag key={t} variant={tagVariant(t)}>{t}</Tag>
              ))}
              <VisibilityTag visibility={recipe.visibility} isDefault={recipe.isDefault} />
            </div>
            <div className="t-display" style={{ fontSize: 52, marginBottom: 12 }}>{recipe.name}</div>
            <div className="t-body muted" style={{ maxWidth: 600 }}>{recipe.description}</div>
            <div className="row gap-2" style={{ marginTop: 20, alignItems: "center" }}>
              <Avatar initials={recipe.author.initials} size="sm" />
              <span className="t-small">oleh <strong>{recipe.author.name}</strong></span>
              <span className="t-small muted">·</span>
              <span className="t-small muted">{recipe.saves} simpan</span>
              <span className="t-small muted">·</span>
              <span className="t-small muted">terakhir brew {recipe.lastBrewed}</span>
            </div>
          </div>
          <div className="row gap-1">
            <button className="icon-btn" title="Favorit"><IconStar size={18} /></button>
            <button className="icon-btn" title="Bagikan"><IconShare size={18} /></button>
            <button className="icon-btn" title="Edit"><IconEdit size={18} /></button>
            <button className="icon-btn" title="Lainnya"><IconKebab size={18} /></button>
          </div>
        </div>

        {/* Params */}
        <Params params={recipe.params} />

        {/* Two-col: timeline + start CTA */}
        <div style={{ display: "grid", gridTemplateColumns: "1.6fr 1fr", gap: 32, marginTop: 32 }}>
          {/* Timeline */}
          <div>
            <div className="row between" style={{ marginBottom: 12, alignItems: "flex-end" }}>
              <div>
                <div className="t-label" style={{ color: "var(--muted)", marginBottom: 6 }}>Alur brewing</div>
                <div className="t-h2">{sessions} sesi · {Math.floor(total / 60)} menit {total % 60} detik total</div>
              </div>
              <button className="btn btn--secondary"><IconEdit size={14} /> Edit alur</button>
            </div>
            <div className="card" style={{ padding: "8px 24px" }}>
              {recipe.timeline.map((step, i) => {
                if (step.kind === "session") sessionIdx++;
                return <StepRow key={i} index={sessionIdx} step={step} />;
              })}
            </div>

            {/* Discussion teaser if group */}
            {recipe.visibility === "group" && (
              <div className="card" style={{ padding: 24, marginTop: 16 }}>
                <div className="row between" style={{ marginBottom: 14 }}>
                  <div className="t-h3">Diskusi grup</div>
                  <span className="t-small muted">3 komentar · 1 disematkan</span>
                </div>
                <div className="row gap-2" style={{ padding: "12px 0", borderTop: "1px solid var(--hairline)" }}>
                  <Avatar initials="RW" size="sm" />
                  <div className="grow">
                    <div className="t-small"><strong>Rizki W.</strong> · 2 hari lalu · <Tag variant="default">DISEMATKAN</Tag></div>
                    <div className="t-small muted" style={{ marginTop: 4 }}>
                      Coba pre-infusion sedikit lebih lama (15 detik) untuk Toraja Sapan — sweetness terangkat lebih bagus.
                    </div>
                  </div>
                </div>
                <div className="row gap-2" style={{ padding: "12px 0", borderTop: "1px solid var(--hairline)" }}>
                  <Avatar initials="MS" size="sm" />
                  <div className="grow">
                    <div className="t-small"><strong>Maya S.</strong> · kemarin</div>
                    <div className="t-small muted" style={{ marginTop: 4 }}>
                      Aku coba di rumah, ratio 1:2.2 untuk Sumatra Lintong rasanya lebih balance. Mau coba update?
                    </div>
                  </div>
                </div>
              </div>
            )}
          </div>

          {/* Start CTA card (dark) */}
          <aside>
            <div className="card card--dark card--hero" style={{ padding: 28, position: "sticky", top: 24 }}>
              <div className="t-label muted-dark" style={{ marginBottom: 12 }}>Siap brewing</div>
              <div className="t-h1" style={{ color: "#fff", marginBottom: 6 }}>
                {Math.floor(total / 60)}<span className="muted-dark" style={{ fontSize: 22 }}>m</span> {total % 60}<span className="muted-dark" style={{ fontSize: 22 }}>s</span>
              </div>
              <div className="t-small muted-dark" style={{ marginBottom: 20 }}>
                Timer akan berjalan otomatis tanpa jeda antar sesi.
              </div>

              <button className="btn btn--primary btn--xl btn--block" onClick={() => onStart(recipe.id)}>
                <IconPlay size={18} /> Mulai Brewing
              </button>

              <div className="row gap-2" style={{ marginTop: 12 }}>
                <button className="btn btn--secondary-dark" style={{ flex: 1 }}>Praktek silent</button>
                <button className="icon-btn icon-btn--dark" title="Reset preferensi"><IconReset size={18} /></button>
              </div>

              <hr className="divider--dark" style={{ margin: "24px 0" }} />

              <div className="t-label muted-dark" style={{ marginBottom: 10 }}>Persiapan</div>
              <ul style={{ margin: 0, padding: 0, listStyle: "none" }}>
                {[
                  `${recipe.params.dose} biji kopi`,
                  `Air ${recipe.params.temp}, ${recipe.params.yield}`,
                  `Grinder set: ${recipe.params.grind}`,
                  `Timer YAVA (otomatis)`,
                ].map((line, i) => (
                  <li key={i} className="row gap-2" style={{ padding: "6px 0", fontSize: 13 }}>
                    <span style={{
                      width: 18, height: 18, borderRadius: "50%",
                      border: "1.5px solid rgba(255,255,255,.30)",
                      display: "inline-flex", alignItems: "center", justifyContent: "center", flex: "none",
                    }}>
                      <span style={{ width: 4, height: 4, borderRadius: "50%", background: "#fff", opacity: .5 }} />
                    </span>
                    <span style={{ color: "rgba(255,255,255,.85)" }}>{line}</span>
                  </li>
                ))}
              </ul>
            </div>
          </aside>
        </div>
      </main>
    </div>
  );
};

window.RecipeDetail = RecipeDetail;
