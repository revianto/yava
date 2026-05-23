// YAVA — mock recipe data
// Hierarchy: Type → Subtype → Recipe with sessions + notes

window.YAVA_DATA = (function () {
  const RECIPES = [
    {
      id: "r-flair-espresso",
      type: "Espresso",
      subtype: "Manual (Flair)",
      name: "Flair Espresso 18g / 36g",
      description:
        "Pre-infusion lembut, ekstraksi 9 bar. Profil seimbang dengan bittersweet finish — cocok untuk single origin Ethiopia.",
      tags: ["ESPRESSO", "MANUAL"],
      visibility: "private",
      isDefault: false,
      params: {
        dose: "18g",
        yield: "36g",
        temp: "94°C",
        grind: "Fine",
        ratio: "1:2",
      },
      timeline: [
        { kind: "session", name: "Pre-infusion", duration: 12, note: "Tekan perlahan — biarkan bubuk basah merata sebelum tekanan penuh." },
        { kind: "session", name: "Extraction", duration: 28, note: "Tekanan stabil 9 bar. Awasi warna mengalir dari emas ke creamy hazel." },
        { kind: "note", content: "Diamkan crema 10 detik sebelum disajikan. Jangan diaduk." },
      ],
      author: { name: "Kamu", initials: "YV" },
      saves: 12,
      lastBrewed: "Kemarin",
    },
    {
      id: "r-v60-light",
      type: "V60",
      subtype: "Regular Drip",
      name: "V60 Light Roast 15g / 250ml",
      description:
        "Pour-over standar untuk light roast. Profil cerah, body ringan, finish clean — disesuaikan untuk bean Kenya AA atau Ethiopia natural.",
      tags: ["V60", "DEFAULT"],
      visibility: "public",
      isDefault: true,
      params: {
        dose: "15g",
        yield: "250ml",
        temp: "92°C",
        grind: "Medium-Fine",
        ratio: "1:16",
      },
      timeline: [
        { kind: "session", name: "Blooming", duration: 45, note: "Tuang 45ml air. Gentle swirl. Nikmati aromanya." },
        { kind: "session", name: "First Pour", duration: 100, note: "Pour ke 150ml total. Lingkaran konsentris dari tengah." },
        { kind: "session", name: "Second Pour", duration: 100, note: "Lanjut ke 250ml. Pour stabil, jangan menyentuh dinding filter." },
      ],
      author: { name: "Sistem YAVA", initials: "YA" },
      saves: 248,
      lastBrewed: "3 hari lalu",
    },
    {
      id: "r-aeropress-inv",
      type: "Aeropress",
      subtype: "Inverted",
      name: "Aeropress Inverted 17g / 220ml",
      description:
        "Metode inverted dengan steep 90 detik. Body lebih kaya dari V60, sweetness terangkat — cocok untuk medium roast.",
      tags: ["AEROPRESS", "INVERTED"],
      visibility: "group",
      isDefault: false,
      params: {
        dose: "17g",
        yield: "220ml",
        temp: "85°C",
        grind: "Medium",
        ratio: "1:13",
      },
      timeline: [
        { kind: "session", name: "Steep", duration: 90, note: "Pour 220ml, aduk 3x dengan stirrer. Diamkan." },
        { kind: "session", name: "Press", duration: 30, note: "Flip dan press perlahan selama 30 detik. Berhenti saat dengar hiss." },
      ],
      author: { name: "Komunitas Senayan", initials: "KS" },
      saves: 67,
      lastBrewed: "Seminggu lalu",
    },
    {
      id: "r-coldbrew",
      type: "Cold Brew",
      subtype: "Slow Drip",
      name: "Cold Brew Slow Drip 50g / 500ml",
      description:
        "Drip selama 8 jam pada suhu ruang. Smooth, low-acid, dark chocolate finish. Bisa disimpan 5 hari di kulkas.",
      tags: ["COLD BREW"],
      visibility: "public",
      isDefault: false,
      params: {
        dose: "50g",
        yield: "500ml",
        temp: "22°C",
        grind: "Coarse",
        ratio: "1:10",
      },
      timeline: [
        { kind: "note", content: "Setup tower drip. Pastikan keran terkalibrasi 1 tetes / 2 detik." },
        { kind: "session", name: "Slow Drip", duration: 60, note: "Demo singkat — versi sebenarnya 8 jam." },
      ],
      author: { name: "Kopi Ruang", initials: "KR" },
      saves: 89,
      lastBrewed: "—",
    },
    {
      id: "r-kopisusu",
      type: "Espresso",
      subtype: "Milk-based",
      name: "Kopi Susu Aren 18g + 150ml",
      description:
        "Double shot dengan gula aren dan susu segar dingin. Manis seimbang, untuk diminum dingin di siang hari.",
      tags: ["ESPRESSO", "GRUP"],
      visibility: "group",
      isDefault: false,
      params: {
        dose: "18g",
        yield: "36g",
        temp: "94°C",
        grind: "Fine",
        ratio: "1:2",
      },
      timeline: [
        { kind: "session", name: "Pre-infusion", duration: 10, note: "Basahi puck merata." },
        { kind: "session", name: "Extraction", duration: 25, note: "9 bar, target 36g." },
        { kind: "note", content: "Tuang 20ml syrup gula aren ke gelas." },
        { kind: "note", content: "Tambahkan ice cubes, lalu susu segar 150ml. Aduk perlahan." },
      ],
      author: { name: "Kamu", initials: "YV" },
      saves: 34,
      lastBrewed: "2 hari lalu",
    },
    {
      id: "r-v60-dark",
      type: "V60",
      subtype: "Regular Drip",
      name: "V60 Dark Roast 18g / 270ml",
      description:
        "Untuk dark roast lokal: Toraja Sapan atau Aceh Gayo. Suhu lebih rendah, ratio sedikit lebih ringan untuk hindari over-extraction.",
      tags: ["V60"],
      visibility: "private",
      isDefault: false,
      params: {
        dose: "18g",
        yield: "270ml",
        temp: "88°C",
        grind: "Medium",
        ratio: "1:15",
      },
      timeline: [
        { kind: "session", name: "Blooming", duration: 40, note: "Tuang 50ml. Swirl ringan." },
        { kind: "session", name: "First Pour", duration: 90, note: "Lanjut ke 160ml." },
        { kind: "session", name: "Second Pour", duration: 80, note: "Tutup di 270ml. Total brew 3:30." },
      ],
      author: { name: "Kamu", initials: "YV" },
      saves: 5,
      lastBrewed: "Hari ini",
    },
  ];

  const TYPES = ["Semua", "Espresso", "V60", "Aeropress", "Cold Brew"];

  const HERO_RECIPE_ID = "r-v60-light";

  return { RECIPES, TYPES, HERO_RECIPE_ID };
})();
