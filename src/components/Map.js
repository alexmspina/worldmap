import React, { useEffect, useRef } from 'react'
import * as L from 'leaflet'
import stylecomponent from '../style/Map/Map.module.css'
import geoJSONmap from '../files/maps/geoJSONmap.json'
import stylemap from '../style/Map/stylemap'
// import targetsFile from '../files/targets/TARGETS_275.json'
import styletarget from '../style/Targets/styletarget'
import { ApolloClient } from 'apollo-client'
import { HttpLink } from 'apollo-link-http'
import { InMemoryCache } from 'apollo-cache-inmemory'
import gql from 'graphql-tag'

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

  useEffect(() => {
    client.query({
      query: gql`
        query {
          targetFeatureCollection {
            type
            features {
              type
              geometry {
                type
                coordinates
              }
              properties {
                shortName
              }
            }
          }
        }
      `
    })
      .then(result => {
        const targetsLayer = L.geoJSON(result.data.targetFeatureCollection, {
          pointToLayer: function (feature, latlng) {
            return L.circleMarker(latlng, styletarget)
          }
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
      query: gql`{
        satellites{
            id,
            latitude,
            longitude,
            velocity,
            altitude,
            mission {
                id,
                config,
                gatewayID,
                gatewayOBAnt,
                gatewayMaxPointingTime,
                beams {
                    id,
                    epcs,
                    targetOBAnt,
                    targetMaxPointingTime,
                    camp,
                    campMode,
                    campGain,
                    ldla,
                    ldlaMode,
                    ldlaFcaGain,
                    ldlaGcaGain,
                    ldlaScaGain
                }
            }
        }
    }`,
      pollInterval: 1000
    })

    satelliteQuery.subscribe({
      next: (result) => console.log(result)
    })
  }, [])

  return (
    <div id='map' className={stylecomponent.map} style={stylecomponent} />
  )
}

export default Map
