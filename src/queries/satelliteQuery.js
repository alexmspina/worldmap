import gql from 'graphql-tag'

const SatelliteQuery = gql`{
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
}`

export default SatelliteQuery
