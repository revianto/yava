import type { Recipe, Group, Discussion, Notification } from '@/types'

export const RECIPES: Recipe[] = [
  {
    id: 'r-flair-espresso',
    type: 'Espresso', subtype: 'Manual (Flair)',
    name: 'Flair Espresso 18g / 36g',
    description: 'Pre-infusion lembut, ekstraksi 9 bar. Profil seimbang dengan bittersweet finish — cocok untuk single origin Ethiopia.',
    tags: ['ESPRESSO', 'MANUAL'],
    visibility: 'private', isDefault: false,
    params: { dose: '18g', yield: '36g', temp: '94°C', grind: 'Fine', ratio: '1:2' },
    timeline: [
      { kind: 'session', name: 'Pre-infusion', duration: 12, note: 'Tekan perlahan — biarkan bubuk basah merata sebelum tekanan penuh.' },
      { kind: 'session', name: 'Extraction', duration: 28, note: 'Tekanan stabil 9 bar. Awasi warna mengalir dari emas ke creamy hazel.' },
      { kind: 'note', content: 'Diamkan crema 10 detik sebelum disajikan. Jangan diaduk.' },
    ],
    author: { name: 'Kamu', initials: 'YV' }, saves: 12, lastBrewed: 'Kemarin',
  },
  {
    id: 'r-v60-light',
    type: 'V60', subtype: 'Regular Drip',
    name: 'V60 Light Roast 15g / 250ml',
    description: 'Pour-over standar untuk light roast. Profil cerah, body ringan, finish clean — disesuaikan untuk bean Kenya AA atau Ethiopia natural.',
    tags: ['V60', 'DEFAULT'],
    visibility: 'public', isDefault: true,
    params: { dose: '15g', yield: '250ml', temp: '92°C', grind: 'Medium-Fine', ratio: '1:16' },
    timeline: [
      { kind: 'session', name: 'Blooming', duration: 45, note: 'Tuang 45ml air. Gentle swirl. Nikmati aromanya.' },
      { kind: 'session', name: 'First Pour', duration: 100, note: 'Pour ke 150ml total. Lingkaran konsentris dari tengah.' },
      { kind: 'session', name: 'Second Pour', duration: 100, note: 'Lanjut ke 250ml. Pour stabil, jangan menyentuh dinding filter.' },
    ],
    author: { name: 'Sistem YAVA', initials: 'YA' }, saves: 248, lastBrewed: '3 hari lalu',
  },
  {
    id: 'r-aeropress-inv',
    type: 'Aeropress', subtype: 'Inverted',
    name: 'Aeropress Inverted 17g / 220ml',
    description: 'Metode inverted dengan steep 90 detik. Body lebih kaya dari V60, sweetness terangkat — cocok untuk medium roast.',
    tags: ['AEROPRESS', 'INVERTED'],
    visibility: 'group', isDefault: false,
    params: { dose: '17g', yield: '220ml', temp: '85°C', grind: 'Medium', ratio: '1:13' },
    timeline: [
      { kind: 'session', name: 'Steep', duration: 90, note: 'Pour 220ml, aduk 3x dengan stirrer. Diamkan.' },
      { kind: 'session', name: 'Press', duration: 30, note: 'Flip dan press perlahan selama 30 detik. Berhenti saat dengar hiss.' },
    ],
    author: { name: 'Komunitas Senayan', initials: 'KS' }, saves: 67, lastBrewed: 'Seminggu lalu',
  },
  {
    id: 'r-coldbrew',
    type: 'Cold Brew', subtype: 'Slow Drip',
    name: 'Cold Brew Slow Drip 50g / 500ml',
    description: 'Drip selama 8 jam pada suhu ruang. Smooth, low-acid, dark chocolate finish. Bisa disimpan 5 hari di kulkas.',
    tags: ['COLD BREW'],
    visibility: 'public', isDefault: false,
    params: { dose: '50g', yield: '500ml', temp: '22°C', grind: 'Coarse', ratio: '1:10' },
    timeline: [
      { kind: 'note', content: 'Setup tower drip. Pastikan keran terkalibrasi 1 tetes / 2 detik.' },
      { kind: 'session', name: 'Slow Drip', duration: 60, note: 'Demo singkat — versi sebenarnya 8 jam.' },
    ],
    author: { name: 'Kopi Ruang', initials: 'KR' }, saves: 89, lastBrewed: '—',
  },
  {
    id: 'r-kopisusu',
    type: 'Espresso', subtype: 'Milk-based',
    name: 'Kopi Susu Aren 18g + 150ml',
    description: 'Double shot dengan gula aren dan susu segar dingin. Manis seimbang, untuk diminum dingin di siang hari.',
    tags: ['ESPRESSO', 'GRUP'],
    visibility: 'group', isDefault: false,
    params: { dose: '18g', yield: '36g', temp: '94°C', grind: 'Fine', ratio: '1:2' },
    timeline: [
      { kind: 'session', name: 'Pre-infusion', duration: 10, note: 'Basahi puck merata.' },
      { kind: 'session', name: 'Extraction', duration: 25, note: '9 bar, target 36g.' },
      { kind: 'note', content: 'Tuang 20ml syrup gula aren ke gelas.' },
      { kind: 'note', content: 'Tambahkan ice cubes, lalu susu segar 150ml. Aduk perlahan.' },
    ],
    author: { name: 'Kamu', initials: 'YV' }, saves: 34, lastBrewed: '2 hari lalu',
  },
  {
    id: 'r-v60-dark',
    type: 'V60', subtype: 'Regular Drip',
    name: 'V60 Dark Roast 18g / 270ml',
    description: 'Untuk dark roast lokal: Toraja Sapan atau Aceh Gayo. Suhu lebih rendah, ratio sedikit lebih ringan untuk hindari over-extraction.',
    tags: ['V60'],
    visibility: 'private', isDefault: false,
    params: { dose: '18g', yield: '270ml', temp: '88°C', grind: 'Medium', ratio: '1:15' },
    timeline: [
      { kind: 'session', name: 'Blooming', duration: 40, note: 'Tuang 50ml. Swirl ringan.' },
      { kind: 'session', name: 'First Pour', duration: 90, note: 'Lanjut ke 160ml.' },
      { kind: 'session', name: 'Second Pour', duration: 80, note: 'Tutup di 270ml. Total brew 3:30.' },
    ],
    author: { name: 'Kamu', initials: 'YV' }, saves: 5, lastBrewed: 'Hari ini',
  },
  {
    id: 'r-chemex-medium',
    type: 'V60', subtype: 'Chemex',
    name: 'Chemex Medium Roast 30g / 500ml',
    description: 'Chemex klasik untuk batch brewing. Filter tebal menghasilkan cup yang sangat bersih. Cocok untuk Papua Wamena atau Flores Bajawa.',
    tags: ['V60', 'CHEMEX'],
    visibility: 'public', isDefault: false,
    params: { dose: '30g', yield: '500ml', temp: '93°C', grind: 'Medium-Coarse', ratio: '1:16' },
    timeline: [
      { kind: 'session', name: 'Bloom', duration: 45, note: 'Tuang 80ml, tunggu gas CO2 keluar.' },
      { kind: 'session', name: 'Pour 1', duration: 90, note: 'Naik ke 250ml, lingkaran dari tengah.' },
      { kind: 'session', name: 'Pour 2', duration: 90, note: 'Tutup di 500ml. Jangan rush.' },
    ],
    author: { name: 'Kopi Ruang', initials: 'KR' }, saves: 42, lastBrewed: '—',
  },
  {
    id: 'r-espresso-default',
    type: 'Espresso', subtype: 'Machine',
    name: 'Espresso Standar 18g / 36g',
    description: 'Resep espresso baseline untuk mesin semi-automatic. Titik awal yang baik untuk dial-in bean apapun.',
    tags: ['ESPRESSO', 'DEFAULT'],
    visibility: 'public', isDefault: true,
    params: { dose: '18g', yield: '36g', temp: '93°C', grind: 'Fine', ratio: '1:2' },
    timeline: [
      { kind: 'session', name: 'Pre-infusion', duration: 8, note: 'Low pressure 3 bar, basahi puck.' },
      { kind: 'session', name: 'Extraction', duration: 25, note: 'Naik ke 9 bar. Target 36g dalam 25–28 detik.' },
    ],
    author: { name: 'Sistem YAVA', initials: 'YA' }, saves: 312, lastBrewed: '—',
  },
  {
    id: 'r-moka-pot',
    type: 'Espresso', subtype: 'Moka Pot',
    name: 'Moka Pot Toraja 20g / 120ml',
    description: 'Moka pot intensitas tinggi. Api kecil, tutup terbuka — untuk dark roast Toraja Sapan yang butuh kontrol suhu ekstra.',
    tags: ['ESPRESSO', 'MOKA'],
    visibility: 'public', isDefault: false,
    params: { dose: '20g', yield: '120ml', temp: '80°C', grind: 'Fine-Medium', ratio: '1:6' },
    timeline: [
      { kind: 'note', content: 'Isi boiler dengan air panas (bukan dingin) untuk hindari over-extraction.' },
      { kind: 'session', name: 'Brew', duration: 90, note: 'Api kecil, tutup dibuka. Angkat saat mulai mendesis.' },
    ],
    author: { name: 'Rina S.', initials: 'RS' }, saves: 78, lastBrewed: '—',
  },
  {
    id: 'r-aeropress-default',
    type: 'Aeropress', subtype: 'Standard',
    name: 'Aeropress Standar 15g / 200ml',
    description: 'Resep Aeropress klasik — standard method, bukan inverted. Hasil konsisten untuk berbagai profil roast.',
    tags: ['AEROPRESS', 'DEFAULT'],
    visibility: 'public', isDefault: true,
    params: { dose: '15g', yield: '200ml', temp: '80°C', grind: 'Medium-Fine', ratio: '1:13' },
    timeline: [
      { kind: 'session', name: 'Bloom', duration: 30, note: 'Tuang 50ml, aduk singkat 3x.' },
      { kind: 'session', name: 'Steep', duration: 60, note: 'Tuang sisa air ke 200ml.' },
      { kind: 'session', name: 'Press', duration: 30, note: 'Press perlahan dalam 20–30 detik.' },
    ],
    author: { name: 'Sistem YAVA', initials: 'YA' }, saves: 189, lastBrewed: '—',
  },
  {
    id: 'r-archived-old',
    type: 'V60', subtype: 'Regular Drip',
    name: 'V60 Percobaan Lama 12g / 200ml',
    description: 'Resep percobaan yang sudah tidak aktif. Ratio terlalu ringan untuk biji yang sekarang dipakai.',
    tags: ['V60'],
    visibility: 'private', isDefault: false, isArchived: true,
    params: { dose: '12g', yield: '200ml', temp: '90°C', grind: 'Medium', ratio: '1:16' },
    timeline: [
      { kind: 'session', name: 'Bloom', duration: 30, note: '' },
      { kind: 'session', name: 'Pour', duration: 120, note: '' },
    ],
    author: { name: 'Kamu', initials: 'YV' }, saves: 0, lastBrewed: '3 bulan lalu',
  },
]

export const GROUPS: Group[] = [
  {
    id: 'g-senayan',
    name: 'Komunitas Senayan',
    description: 'Komunitas kopi rumahan di sekitar Senayan & Kebayoran. Sharing resep, tips dial-in, dan occasional cupping session.',
    inviteCode: 'SNYN2026',
    myRole: 'member',
    createdAt: 'Jan 2026',
    members: [
      { id: 'm-1', name: 'Rizki W.', initials: 'RW', role: 'founder', joinedAt: 'Jan 2026' },
      { id: 'm-2', name: 'Adi P.', initials: 'AP', role: 'admin', joinedAt: 'Jan 2026' },
      { id: 'm-3', name: 'Maya S.', initials: 'MS', role: 'member', joinedAt: 'Feb 2026' },
      { id: 'm-4', name: 'Dian K.', initials: 'DK', role: 'member', joinedAt: 'Feb 2026' },
      { id: 'm-5', name: 'Budi H.', initials: 'BH', role: 'member', joinedAt: 'Mar 2026' },
      { id: 'm-me', name: 'Nadira (Kamu)', initials: 'ND', role: 'member', joinedAt: 'Mar 2026' },
    ],
    recipes: [
      { id: 'gr-1', recipeId: 'r-aeropress-inv', recipeName: 'Aeropress Inverted 17g / 220ml', recipeType: 'Aeropress', recipeSubtype: 'Inverted', submittedBy: 'Kamu', submittedByInitials: 'ND', submittedAt: '2 minggu lalu', status: 'active' },
      { id: 'gr-2', recipeId: 'r-kopisusu', recipeName: 'Kopi Susu Aren 18g + 150ml', recipeType: 'Espresso', recipeSubtype: 'Milk-based', submittedBy: 'Maya S.', submittedByInitials: 'MS', submittedAt: '3 minggu lalu', status: 'active' },
      { id: 'gr-3', recipeId: 'r-v60-dark', recipeName: 'V60 Dark Roast 18g / 270ml', recipeType: 'V60', recipeSubtype: 'Regular Drip', submittedBy: 'Adi P.', submittedByInitials: 'AP', submittedAt: '1 bulan lalu', status: 'active' },
      { id: 'gr-4', recipeId: 'r-coldbrew', recipeName: 'Cold Brew Slow Drip 50g / 500ml', recipeType: 'Cold Brew', recipeSubtype: 'Slow Drip', submittedBy: 'Budi H.', submittedByInitials: 'BH', submittedAt: '5 hari lalu', status: 'pending' },
      { id: 'gr-5', recipeId: 'r-moka-pot', recipeName: 'Moka Pot Toraja 20g / 120ml', recipeType: 'Espresso', recipeSubtype: 'Moka Pot', submittedBy: 'Dian K.', submittedByInitials: 'DK', submittedAt: '2 hari lalu', status: 'pending' },
      { id: 'gr-6', recipeId: 'r-flair-espresso', recipeName: 'Flair Espresso 18g / 36g', recipeType: 'Espresso', recipeSubtype: 'Manual (Flair)', submittedBy: 'Rizki W.', submittedByInitials: 'RW', submittedAt: '1 bulan lalu', status: 'rejected', rejectionReason: 'Ratio kurang sesuai untuk biji Sumatra — coba 1:2.2 dulu.' },
    ],
  },
  {
    id: 'g-barista',
    name: 'YAVA Barista Club',
    description: 'Grup kecil untuk eksperimen resep advanced. Filter, espresso, dan cold brew specialty.',
    inviteCode: 'YBC2026',
    myRole: 'founder',
    createdAt: 'Apr 2026',
    members: [
      { id: 'b-me', name: 'Nadira (Kamu)', initials: 'ND', role: 'founder', joinedAt: 'Apr 2026' },
      { id: 'b-2', name: 'Sinta R.', initials: 'SR', role: 'admin', joinedAt: 'Apr 2026' },
      { id: 'b-3', name: 'Hendra P.', initials: 'HP', role: 'member', joinedAt: 'Mei 2026' },
      { id: 'b-4', name: 'Layla N.', initials: 'LN', role: 'member', joinedAt: 'Mei 2026' },
    ],
    recipes: [
      { id: 'br-1', recipeId: 'r-v60-light', recipeName: 'V60 Light Roast 15g / 250ml', recipeType: 'V60', recipeSubtype: 'Regular Drip', submittedBy: 'Kamu', submittedByInitials: 'ND', submittedAt: '1 minggu lalu', status: 'active' },
      { id: 'br-2', recipeId: 'r-chemex-medium', recipeName: 'Chemex Medium Roast 30g / 500ml', recipeType: 'V60', recipeSubtype: 'Chemex', submittedBy: 'Sinta R.', submittedByInitials: 'SR', submittedAt: '3 hari lalu', status: 'pending' },
    ],
  },
]

export const DISCUSSIONS: Discussion[] = [
  {
    id: 'd-1', recipeId: 'r-aeropress-inv',
    authorName: 'Rizki W.', authorInitials: 'RW',
    content: 'Coba pre-infusion sedikit lebih lama (15 detik) untuk Toraja Sapan — sweetness terangkat lebih bagus.',
    createdAt: '2 hari lalu', pinned: true,
    replies: [
      { id: 'd-1-r1', authorName: 'Maya S.', authorInitials: 'MS', content: 'Setuju! Aku coba dan beda banget hasilnya. Pre-infusion 15-20 detik jadi sweet spot.', createdAt: 'kemarin' },
    ],
  },
  {
    id: 'd-2', recipeId: 'r-aeropress-inv',
    authorName: 'Maya S.', authorInitials: 'MS',
    content: 'Aku coba di rumah, ratio 1:2.2 untuk Sumatra Lintong rasanya lebih balance. Mau coba update?',
    createdAt: 'kemarin', pinned: false,
    replies: [],
  },
  {
    id: 'd-3', recipeId: 'r-kopisusu',
    authorName: 'Dito A.', authorInitials: 'DA',
    content: 'Syrup gula aren berapa ml tepatnya? 20ml terasa kurang manis buat aku.',
    createdAt: '3 hari lalu', pinned: false,
    replies: [
      { id: 'd-3-r1', authorName: 'Kamu', authorInitials: 'ND', content: 'Aku biasanya 25ml kalau pakai bean yang lebih asam seperti Kenya. Coba adjust sesuai selera.', createdAt: '3 hari lalu' },
      { id: 'd-3-r2', authorName: 'Dito A.', authorInitials: 'DA', content: 'Thanks! 25ml pas banget.', createdAt: '2 hari lalu' },
    ],
  },
]

export const NOTIFICATIONS: Notification[] = [
  { id: 'n-1', type: 'approved', title: 'Resep disetujui', body: 'V60 Light Roast 15g / 250ml kamu disetujui di Komunitas Senayan.', link: '/groups/g-senayan', read: false, createdAt: '1 jam lalu' },
  { id: 'n-2', type: 'reply', title: 'Balas diskusi', body: 'Maya S. membalas komentar kamu di Aeropress Inverted 17g / 220ml.', link: '/recipes/r-aeropress-inv', read: false, createdAt: '3 jam lalu' },
  { id: 'n-3', type: 'joined', title: 'Anggota baru', body: 'Layla N. bergabung ke YAVA Barista Club.', link: '/groups/g-barista', read: false, createdAt: 'kemarin' },
  { id: 'n-4', type: 'rejected', title: 'Resep ditolak', body: 'Cold Brew Slow Drip ditolak di Komunitas Senayan dengan catatan: perlu revisi parameter grind.', link: '/groups/g-senayan', read: true, createdAt: '2 hari lalu' },
  { id: 'n-5', type: 'approved', title: 'Resep disetujui', body: 'Kopi Susu Aren 18g + 150ml kamu disetujui di Komunitas Senayan.', link: '/groups/g-senayan', read: true, createdAt: '3 hari lalu' },
]

export const TYPES = ['Semua', 'Espresso', 'V60', 'Aeropress', 'Cold Brew']
export const HERO_RECIPE_ID = 'r-v60-light'

export const totalDuration = (r: Recipe) =>
  r.timeline.filter((s): s is import('@/types').RecipeSession => s.kind === 'session').reduce((a, s) => a + s.duration, 0)

export const sessionCount = (r: Recipe) =>
  r.timeline.filter((s) => s.kind === 'session').length

export const tagVariant = (label: string): string => {
  const k = label.toUpperCase()
  if (k === 'ESPRESSO') return 'espresso'
  if (k === 'V60') return 'v60'
  if (k.includes('COLD')) return 'cold'
  if (k === 'GRUP') return 'group'
  return 'default'
}

export const fmt = (seconds: number): string => {
  const s = Math.max(0, Math.round(seconds))
  const m = Math.floor(s / 60)
  const r = s % 60
  return `${String(m).padStart(2, '0')}:${String(r).padStart(2, '0')}`
}
