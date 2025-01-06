package godio

type Interval int

const (
	PerfectUnison    Interval = 0
	MinorSecond      Interval = 1
	MajorSecond      Interval = 2
	MinorThird       Interval = 3
	MajorThird       Interval = 4
	PerfectFourth    Interval = 5
	DiminishedFifth  Interval = 6
	PerfectFifth     Interval = 7
	MinorSixth       Interval = 8
	MajorSixth       Interval = 9
	MinorSeventh     Interval = 10
	MajorSeventh     Interval = 11
	PerfectOctave    Interval = 12
	MinorNinth       Interval = 13
	MajorNinth       Interval = 14
	MinorTenth       Interval = 15
	MajorTenth       Interval = 16
	PerfectEleventh  Interval = 17
	AugmentedEleven  Interval = 18
	MinorThirteenth  Interval = 20
	MajorThirteenth  Interval = 21
	MinorFourteenth  Interval = 22
	MajorFourteenth  Interval = 23
	PerfectFifteenth Interval = 24
)

var IntervalNames = map[string]Interval{
	"P1": PerfectUnison,
	"m2": MinorSecond,
	"M2": MajorSecond,
	"m3": MinorThird,
	"M3": MajorThird,
	"P4": PerfectFourth,
	"A4": DiminishedFifth,
	"P5": PerfectFifth,
	"m6": MinorSixth,
	"M6": MajorSixth,
	"m7": MinorSeventh,
	"M7": MajorSeventh,
	"P8": PerfectOctave,
}
