export type JourneyYear = `${number}`
export type JourneyMonth = 'Jan' | 'Feb' | 'Mar' | 'Apr' | 'May' | 'Jun' | 'Jul' | 'Aug' | 'Sep' | 'Oct' | 'Nov' | 'Dec'

export interface JourneyPoint {
  year: JourneyYear
  month: JourneyMonth
}

export interface Journey {
  id: number
  year: JourneyYear
  month: JourneyMonth
  topic: string
  place: string
  date: string
  thumbnail: string | null
}

export interface JourneyTimelineItem {
  year: JourneyYear
  months: JourneyMonth[]
}
