package rex

// Pipe allows us to build data pipelines
type Pipe struct {
	data chan Record
}

type RecordProducer chan Record
type RecordProcessor func(Record) Record
type RecordCollector func(chan Record)
type RecordFilter func(Record) bool

// Unue ("first") start with a data producer.
func Unue(producer RecordProducer) Pipe {
	return Pipe{
		data: producer,
	}
}

// Kaj ("and") process each input record with the processor.
func (p Pipe) Kaj(process RecordProcessor) Pipe {

	c := make(chan Record)
	go func() {
		defer close(c)

		for r := range p.data {
			c <- process(r)
		}
	}()

	return Pipe{
		data: c,
	}
}

// Inkluzivi ("include") checks a boolean function against each record to see if
// it should be included.
func (p Pipe) Inkluzivi(f RecordFilter) Pipe {
	return p.Kaj(func(r Record) Record {
		r.dropped = r.dropped || !f(r)
		return r
	})
}

// Ekskludi ("exclude") checks a boolean function against each record to see if
// it should be included.
func (p Pipe) Ekskludi(f RecordFilter) Pipe {
	return p.Kaj(func(r Record) Record {
		r.dropped = r.dropped || f(r)
		return r
	})
}

func drop(in chan Record) chan Record {

	c := make(chan Record)

	go func() {
		defer close(c)

		for r := range in {
			if r.dropped {
				continue
			}

			c <- r
		}
	}()

	return c
}

// Kolekti ("collect") invoke a method that will consume all the data.
func (p Pipe) Kolekti(collector RecordCollector) {
	collector(drop(p.data))
}

// Fine ("finally") is the final step of a pipeline. Consume all the results.
func (p Pipe) Fine() {
	for range p.data {
		// nenio ("nothing")
	}
}
