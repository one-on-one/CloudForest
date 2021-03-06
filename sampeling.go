package CloudForest

import (
	"math/rand"
)

type Bagger interface {
	Sample(samples *[]int, n int)
}

//BalancedSampler provides for random sampelign of integers (usually case indexes)
//in a way that ensures a balanced presence of classes.
type BalancedSampler struct {
	Cases [][]int
}

//NeaBalancedSampler initalizes a balanced sampler that will evenly balance cases
//between the classes present in the provided DesnseeCatFeature.
func NewBalancedSampler(catf *DenseCatFeature) (s *BalancedSampler) {
	s = &BalancedSampler{make([][]int, catf.NCats())}

	for i := 0; i < catf.NCats(); i++ {
		s.Cases = append(s.Cases, make([]int, catf.Length()))
	}

	for i, v := range catf.CatData {
		if !catf.IsMissing(i) {
			s.Cases[v] = append(s.Cases[v], i)
		}
	}
	return
}

//Sample samples n integers in a balnced-with-replacment fashion into the provided array.
func (s *BalancedSampler) Sample(samples *[]int, n int) {
	(*samples) = (*samples)[0:0]
	nCases := len(s.Cases)
	c := 0
	for i := 0; i < n; i++ {
		c = rand.Intn(nCases)
		(*samples) = append((*samples), s.Cases[c][rand.Intn(len(s.Cases[c]))])
	}

}

//A secondary balanced sampler roughly balances the target feature within the classes of another catagorical
//feature while roughly preserving the origional rate of the secondary feature.
type SecondaryBalancedSampler struct {
	Total    int
	Counts   []int
	Samplers [][][]int
}

func NewSecondaryBalancedSampler(target *DenseCatFeature, balanceby *DenseCatFeature) (s *SecondaryBalancedSampler) {
	nSecondaryCats := balanceby.NCats()
	s = &SecondaryBalancedSampler{0, make([]int, nSecondaryCats, nSecondaryCats), make([][][]int, nSecondaryCats)}

	for i := 0; i < nSecondaryCats; i++ {
		s.Samplers = append(s.Samplers, make([][]int, target.NCats()))
		for j := 0; j < target.NCats(); j++ {
			s.Samplers[i] = append(s.Samplers[i], make([]int, target.Length()))
		}

	}

	for i := 0; i < target.Length(); i++ {
		if !target.IsMissing(i) && !balanceby.IsMissing(i) {
			s.Total += 1
			balanceCat := balanceby.Geti(i)
			targetCat := target.Geti(i)
			s.Counts[balanceCat] += 1
			s.Samplers[balanceCat][targetCat] = append(s.Samplers[balanceCat][targetCat], i)
		}
	}
	return

}

func (s *SecondaryBalancedSampler) Sample(samples *[]int, n int) {
	(*samples) = (*samples)[0:0]

	b := 0
	c := 0
	for i := 0; i < n; i++ {
		b = rand.Intn(s.Total)
		for j, v := range s.Counts {
			b = b - v
			if b < 0 || j == (len(s.Counts)-1) {
				b = j
				break
			}
		}
		nCases := len(s.Samplers[b])
		c = rand.Intn(nCases)
		(*samples) = append((*samples), s.Samplers[b][c][rand.Intn(len(s.Samplers[b][c]))])
	}

}

/*
SampleFirstN ensures that the first n entries in the supplied
deck are randomly drawn from all entries without replacement for use in selecting candidate
features to split on. It accepts a pointer to the deck so that it can be used repeatedly on
the same deck avoiding reallocations.
*/
func SampleFirstN(deck *[]int, n int) {
	cards := *deck
	length := len(cards)
	old := 0
	randi := 0
	for i := 0; i < n; i++ {
		old = cards[i]
		randi = i + rand.Intn(length-i)
		cards[i] = cards[randi]
		cards[randi] = old

	}
}

/*
SampleWithReplacment samples nSamples random draws from [0,totalCases) with replacement
for use in selecting cases to grow a tree from.
*/
func SampleWithReplacment(nSamples int, totalCases int) (cases []int) {
	cases = make([]int, 0, nSamples)
	for i := 0; i < nSamples; i++ {
		cases = append(cases, rand.Intn(totalCases))
	}
	return
}
