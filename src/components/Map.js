import React, { useEffect, useRef } from 'react'
import L from 'leaflet'
import style from '../style/Map/Map.module.css'
import mapSVG from '../images/maps/map02.svg'
// import ZoomControl from 'react-leaflet/lib/ZoomControl'
// import AttributionControl from 'react-leaflet/lib/AttributionControl'

const imageBounds = L.latLngBounds([[90, -180], [-90, 180]])

function Map ({ markerPosition }) {
  // create map
  const mapRef = useRef(null)
  const overlay = L.imageOverlay(mapSVG, imageBounds)
  useEffect(() => {
    mapRef.current = L.map('map', {
      zoom: 1,
      minZoom: 1,
      zoomControl: false,
      attributionControl: false,
      crs: L.CRS.Simple
    })

    mapRef.current.addLayer(overlay)
    mapRef.current.fitBounds(imageBounds)
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
