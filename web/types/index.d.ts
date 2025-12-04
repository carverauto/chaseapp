import { GeoPoint } from 'firebase/firestore'

export interface FirehoseData {
    actor: string
    created_at: Date
    eventType: string
    foreign_id?: string
    id: string
    object: string
    origin?: any
    target?: any
    time: Date
    verb: string
    payload: object
}

export interface State {
    firehoseData: FirehoseData
    user: AuthUser
    chase: Chase
    chases: Chase[]
    chaseId: string
    lastUpdateDate: number
    mapCoords: number[]
    count: number
    authUser: AuthUser
    messagingToken: string
    streamToken: string
}

export interface Marker {
    name: string
    desc?: string
    type: string
    location: number[]
    urls?: string[]
}

export interface AuthUser {
    uid: string
    displayName?: string
    userName?: string
    email?: string
    emailVerified?: boolean
    photoURL?: string
    isAdmin?: object
    favorites?: string[]
}

export interface FirehoseEventPayload {
    name: string
    urls: string[]
}

export interface FirehoseEvent {
    id?: string
    createdAt?: Date
    eventType?: string
    payload?: FirehoseEventPayload
}
export interface LiveATC {
    name: string
    url: string
}

export interface Airport {
    airport: string
    city: string
    state: string
    iata: string
    icao: string
    liveatc: LiveATC[]
    location: GeoPoint
}

export enum AirshipType {
    heli,
    plane,
    ac130
}

export enum AirshipGroup {
    leo,
    media,
    fire,
    rescue,
    noaa
}

export interface Airship {
    group: AirshipGroup
    imageUrl: string
    tailno: string
    type: AirshipType
}

export interface BoFbird {
    tailno: string
    distance: number
}

export interface BoFPayload {
    leo: BoFbird[]
}

export interface BoF {
    createdAt: Date
    eventType: string
    payload: BoFPayload[]
}

export interface Sentiment {
    magnitude: number
    score: number
}

export interface Wheels {
    W1: string
    W2: string
    W3: string
    W4: string
}

export interface Chase {
    ID: string
    ImageURL?: string
    Networks?: Network[]
    Desc?: string
    Name?: string
    // CreatedAt?: FirebaseFirestore.Timestamp
    CreatedAt?: Date
    EndedAt?: Date
    Votes?: number
    Live?: boolean
    sentiment?: Sentiment
    Type?: string
    Wheels?: Wheels
}

export interface Rocket {
    Owner: string
    Reused: number
}

export interface Coordinates {
    lat: number
    lng: number
}

export interface Streams {
    Tier: string
    URL: string
}

export interface Network {
    Logo: string
    Name: string
    Other: string
    Tier: number
    URL: string
    Streams: Streams[]
}

export interface Vehicle {
    CompanyID: number
    ID: number
    Name: string
    Slug: string
}

export interface Provider {
    ID: number
    Name: string
    Slug: string
}

export interface Launch {
    id: string
    ImageURL?: string
    Live: boolean
    Votes: number
    Desc: string
    Vehicle: Vehicle
    LaunchDesc: string
    Provider: Provider
    Name: string
    Networks: Network[]
    Rocket: Rocket
    CreatedAt: Date
    LaunchDate: Date
    ActualLaunchDate: Date
    Status: string
    Coordinates: Coordinates
}

export interface ADSB {
    alt_baro: number
    alt_geom: number
    alert1: number
    emergency: number
    flight: number
    geohash: string
    group: string
    gs: number
    hex: string
    imageUrl: string
    lat: number
    lon: number
    messages: number
    rssi: number
    seen: number
    seen_pos: number
    squawk: number
    t: string
    tailno: string
    track: number
    type: string
    updated: Date
    version: number
    sprite?: string
}

export interface GeojsonGeometry {
    type: string
    coordinates: number[]
}

export interface GeojsonProperties {
    title: string
    type?: string
    imageUrl?: string
    group?: string
    dbscan?: string
    cluster?: number
}

export interface GeojsonFeature {
    type: string
    geometry: GeojsonGeometry
    properties: GeojsonProperties
}

export interface MyGeojsonData {
    type: 'FeatureCollection'
    features: GeojsonFeature[]
}

export interface MyGeojson {
    type: 'geojson'
    data: MyGeojsonData
}

export interface WeatherStationBenchmark {
    self: string
}

export interface WeatherStationSensors {
    self: string
}
export interface WeatherStation {
    stormsurge: boolean
    id: number
    lat: number
    lng: number
    name: string
    sensors: WeatherStationSensors
}

export interface weather {
    stations: WeatherStation
}

export interface Cluster {
    clusterId: number
    group: string
    tailNo: string
}

export interface AisBoat {
    a: number
    b: number
    c: number
    d: number
    callsign: string
    cog: number
    dest: string
    draught: number
    eta: string
    heading: number
    imo: number
    latitude: number
    longitude: number
    mmsi: string
    name: string
    nav_sat: number
    rot: number
    sog: number
    type: number
}
