import L from 'leaflet'

var iconSettings = {
  mapIconUrl: `
          <svg id='satellite' viewBox='0 0 131.3 113.5' className='footer__form__button__svg'>
            <path className='footer__button__svg__path' d='M58.5,76.3l-10.3-8.2l25.2-32l10.3,8.2L58.5,76.3z M99,65.7l-17.5,3.4L82,87.5l31.3,25l17.1-21.8L99,65.7z
                  M32.3,47.8L32.3,47.8l17.5-3.4L49.4,26L18.1,1L1,22.8L32.3,47.8z M81.5,69.1l-10.7-8.3 M49.7,44.2l9.9,7.7' stroke='#48aaf2' fill='#994ae8' stroke-width='5' />
          </svg>
        `
}

// icon normal state
var satelliteDivIcon = L.divIcon({
  className: 'leaflet-data-marker',
  html: L.Util.template(iconSettings.mapIconUrl, iconSettings), // .replace('#','%23'),
  iconAnchor: [30, 10],
  iconSize: [40, 40],
  popupAnchor: [0, -28]
})

export default satelliteDivIcon
