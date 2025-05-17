package graphqlcontroller

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/DKhorkov/libs/logging"

	"github.com/DKhorkov/hmtm-bff/internal/config"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
)

func TestNewResolver(t *testing.T) {
	testCases := []struct {
		name          string
		useCases      interfaces.UseCases
		logger        logging.Logger
		cookiesConfig config.CookiesConfig
		expected      *Resolver
	}{
		{
			name:     "success",
			expected: &Resolver{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := NewResolver(tc.useCases, tc.logger, tc.cookiesConfig)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestResolver_Toy(t *testing.T) {
	testCases := []struct {
		name     string
		resolver *Resolver
		expected *toyResolver
	}{
		{
			name:     "success",
			resolver: &Resolver{},
			expected: &toyResolver{Resolver: &Resolver{}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.resolver.Toy()
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestResolver_Ticket(t *testing.T) {
	testCases := []struct {
		name     string
		resolver *Resolver
		expected *ticketResolver
	}{
		{
			name:     "success",
			resolver: &Resolver{},
			expected: &ticketResolver{Resolver: &Resolver{}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.resolver.Ticket()
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestResolver_Respond(t *testing.T) {
	testCases := []struct {
		name     string
		resolver *Resolver
		expected *respondResolver
	}{
		{
			name:     "success",
			resolver: &Resolver{},
			expected: &respondResolver{Resolver: &Resolver{}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.resolver.Respond()
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestResolver_Query(t *testing.T) {
	testCases := []struct {
		name     string
		resolver *Resolver
		expected *queryResolver
	}{
		{
			name:     "success",
			resolver: &Resolver{},
			expected: &queryResolver{Resolver: &Resolver{}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.resolver.Query()
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestResolver_Mutation(t *testing.T) {
	testCases := []struct {
		name     string
		resolver *Resolver
		expected *mutationResolver
	}{
		{
			name:     "success",
			resolver: &Resolver{},
			expected: &mutationResolver{Resolver: &Resolver{}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.resolver.Mutation()
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestResolver_Master(t *testing.T) {
	testCases := []struct {
		name     string
		resolver *Resolver
		expected *masterResolver
	}{
		{
			name:     "success",
			resolver: &Resolver{},
			expected: &masterResolver{Resolver: &Resolver{}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.resolver.Master()
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestResolver_Email(t *testing.T) {
	testCases := []struct {
		name     string
		resolver *Resolver
		expected *emailResolver
	}{
		{
			name:     "success",
			resolver: &Resolver{},
			expected: &emailResolver{Resolver: &Resolver{}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.resolver.Email()
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestResolver_Pagination(t *testing.T) {
	testCases := []struct {
		name     string
		resolver *Resolver
		expected *paginationResolver
	}{
		{
			name:     "success",
			resolver: &Resolver{},
			expected: &paginationResolver{Resolver: &Resolver{}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.resolver.Pagination()
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestResolver_ToysFilters(t *testing.T) {
	testCases := []struct {
		name     string
		resolver *Resolver
		expected *toysFiltersResolver
	}{
		{
			name:     "success",
			resolver: &Resolver{},
			expected: &toysFiltersResolver{Resolver: &Resolver{}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.resolver.ToysFilters()
			require.Equal(t, tc.expected, actual)
		})
	}
}
