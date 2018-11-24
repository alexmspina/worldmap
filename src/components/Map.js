import React, { useEffect, useRef } from 'react'
import L from 'leaflet'
import style from '../style/Map/Map.module.css'
import mapSVG from '../images/maps/map02.svg'
// import ZoomControl from 'react-leaflet/lib/ZoomControl'
// import AttributionControl from 'react-leaflet/lib/AttributionControl'

const mapBounds = L.latLngBounds([[45, -180], [-45, 180]])

const imageBounds = L.latLngBounds([[-120, -200], [120, 200]])

function Map ({ markerPosition }) {
  // create map
  const mapRef = useRef(null)
  const overlay = L.imageOverlay(mapSVG, imageBounds, {
    opacity: 0.5
  })
  useEffect(() => {
    mapRef.current = L.map('map', {
      zoomControl: false,
      attributionControl: false,
      zoomSnap: 0.25,
      zoom: 2.25,
      layers: [
        L.tileLayer('http://{s}.tile.osm.org/{z}/{x}/{y}.png', {
          attribution:
            '&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'
        })
      ]
    })

    mapRef.current.addLayer(overlay)
    mapRef.current.fitBounds(mapBounds)
  }, [])

  // add marker
  const markerRef = useRef(null)
  useEffect(() => {
    if (markerRef.current) {
      markerRef.current.setLatLng(markerPosition)
    } else {
      markerRef.current = L.marker(markerPosition).addTo(mapRef.current)
    }
  },
  [markerPosition]
  )

  return <div id='map' className={style.map} style={style} />
}

export default Map

// L.tileLayer('http://{s}.tile.osm.org/{z}/{x}/{y}.png', {
//           attribution:
//             '&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'
//         })
