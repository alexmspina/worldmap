import React, { useEffect, useRef } from 'react'
import L from 'leaflet'

// Map outline in form of geojson points
import geoJSONmap from '../files/maps/geoJSONmap.json'

// Style
import 'leaflet/dist/leaflet.css'
import stylecomponent from '../style/Map/Map.module.css'
import stylemap from '../style/Map/stylemap'
import styletarget from '../style/Targets/styletarget'
// import stylesatellite from '../style/Satellites/stylesatellites'
import satelliteDivIcon from './../images/icons/satelliteSVG'
import gatewayDivIcon from './../images/icons/gatewaySVG'

// Apollo
import { ApolloClient } from 'apollo-client'
import { HttpLink } from 'apollo-link-http'
import { InMemoryCache } from 'apollo-cache-inmemory'

// Queries
import SatelliteQuery from './../queries/satelliteQuery'
import TargetsQuery from './../queries/targetsQuery'

function Map () {
  const client = new ApolloClient({
    uri: 'http://localhost:3000/graphql',
    link: new HttpLink(),
    cache: new InMemoryCache({
      addTypename: false
    })
  })

  // create map reference
  const mapRef = useRef(null)

  // create map from geoJSON layer
  const maplayer = L.geoJSON(geoJSONmap, {
    style: stylemap
  })

  let newSatelliteLayer = L.geoJSON([], {
    pointToLayer: function (feature, latlng) {
      return L.marker(latlng, { icon: satelliteDivIcon })
    }
  })

  useEffect(() => {
    client.query({
      query: TargetsQuery
    })
      .then(result => {
        const targetsLayer = L.geoJSON(result.data.targetFeatureCollection, {
          pointToLayer: function (feature, latlng) {
            console.log(feature.properties)
            if (feature.properties.gatewayFlag === 'G') {
              const icon = gatewayDivIcon
              return L.marker(latlng, { icon: icon })
            }
            return L.circleMarker(latlng, styletarget)
          }
        }).bindPopup(function (layer) {
          return (
            `
            <div>
              <h1>
                ${layer.feature.properties.shortName}
              </h1>
            `
          )
        })

        mapRef.current = L.map('map', {
          zoomControl: false,
          attributionControl: false,
          zoomSnap: 0.25,
          minZoom: 2.5,
          boxZoom: true
        })
        mapRef.current.addLayer(maplayer)
        mapRef.current.addLayer(targetsLayer)
        mapRef.current.setView([0, 0], 2.5)
        mapRef.current.setMaxBounds(mapRef.current.getBounds())
      })
      .catch(error => console.error(error))
  }, [])

  useEffect(() => {
    const satelliteQuery = client.watchQuery({
      query: SatelliteQuery,
      pollInterval: 1000
    })

    const plotSats = (satellites, layer) => {
      if (mapRef.current !== null) {
        layer.clearLayers()
        layer.addData(satellites)
        layer.bindPopup(function (layer) {
          return (
            `
            <div>
              <h1>
                ${layer.feature.properties.id}
              </h1>
              <div>
                <h2>
                  Longitude:
                </h2>
                ${Math.round(layer.feature.geometry.coordinates[0])}
              </div>
              <div>
                <h2>
                  Missions
                </h2>
                <div>
                  ${layer.feature.properties.mission.map(mission => mission.beams.map(beam => beam.id))}
                </div>
              </div>
            </div>
            `
          )
        })
        mapRef.current.addLayer(layer)
      }
    }

    satelliteQuery.subscribe({
      next: (result) => plotSats(result.data.satelliteFeatureCollection, newSatelliteLayer)
    })
  }, [])

  return (
    <div id='map' className={stylecomponent.map} style={stylecomponent} />
  )
}

export default Map
