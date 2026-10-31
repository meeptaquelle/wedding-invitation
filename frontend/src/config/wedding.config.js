// CONTENT EDIT:
// All wedding information lives here. Templates never duplicate this data.
export const weddingConfig = {
  couple: {
    partnerOne: 'Prabowo',
    partnerTwo: 'Teddy'
  },

  // CONTENT EDIT: Event information
  event: {
    date: '2027-01-01',
    time: '10:00',
    timezone: 'Asia/Makassar',
    venue: 'The Grand Ballroom',
    address: 'Jl. Contoh Raya No. 123, Makassar',
    mapsUrl: 'https://maps.google.com/?q=Grand+Ballroom+Makassar'
  },

  // DESIGN EDIT: Image assets
  images: {
    hero: '/images/examplephoto.jpg',
    gallery: [
      '/images/examplepotrait.jpg',
      '/images/example-gallery-2.svg',
      '/images/example-gallery-3.svg'
    ]
  },

  // DESIGN EDIT: Optional sections
  sections: {
    introduction: true,
    gallery: true,
    countdown: true,
    guestbook: true
  },

  // CONTENT EDIT: Introduction text
  introduction:
    'With joy in our hearts, we invite you to share in the celebration of our wedding. Your presence and blessings mean the world to us.',

  // CONTENT EDIT: Footer text
  footer: {
    closing: 'With love,',
    signature: 'Prabowo & Teddy',
    credit: 'Made with love'
  }
}
