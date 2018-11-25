import React, { useEffect, useRef } from 'react'
import * as L from 'leaflet'
import style from '../style/Map/Map.module.css'
import mapSVG from '../images/maps/map02.svg'

const imageBounds = L.latLngBounds([[-70, -186.5], [82.5, 215.5]])

const point = {
  'type': 'FeatureCollection',
  'features': [{
    'type': 'Feature',
    'geometry': {
      'type': 'Point',
      'coordinates': [-5.65, 35.95]
    },
    'properties': {
      'prop0': 'value0'
    }
  }]
}

const geoStyle = {
  radius: 8,
  fillColor: '#ff7800',
  color: '#000',
  weight: 1,
  opacity: 1,
  fillOpacity: 0.8
}

function Map ({ markerPosition }) {
  // create map
  const mapRef = useRef(null)
  const overlay = L.imageOverlay(mapSVG, imageBounds, {
    opacity: 0.5
  })

  const geolayer = L.geoJSON(point, {
    pointToLayer: function (feature, latlng) {
      return L.circleMarker(latlng, geoStyle)
    }
  })

  useEffect(() => {
    mapRef.current = L.map('map', {
      zoomControl: false,
      attributionControl: false,
      zoomSnap: 0.25,
      minZoom: 2.25,
      boxZoom: true
    })

    mapRef.current.addLayer(overlay)
    mapRef.current.addLayer(geolayer)
    mapRef.current.setView([10, 10], 2.25)
  }, [])

  return <div id='map' className={style.map} style={style} />
}

export default Map

// L.tileLayer('http://{s}.tile.osm.org/{z}/{x}/{y}.png', {
//           attribution:
//             '&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'
//         })
