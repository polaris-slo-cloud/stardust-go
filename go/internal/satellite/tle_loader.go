package satellite

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/keniack/stardustGo/configs"
	"github.com/keniack/stardustGo/internal/links"
	"github.com/keniack/stardustGo/pkg/types"
)

const (
	dataSourceType = "tle"
	errCannotParse = "cannot parse tle data source"
)

// TleLoader reads and parses satellites from a TLE (Two-Line Element) data source.
type TleLoader struct {
	config           configs.InterSatelliteLinkConfig
	satelliteBuilder *SatelliteBuilder
}

// NewTleLoader creates a new TleLoader instance.
func NewTleLoader(config configs.InterSatelliteLinkConfig, builder *SatelliteBuilder) *TleLoader {
	return &TleLoader{
		config:           config,
		satelliteBuilder: builder,
	}
}

// Load parses the TLE stream into Satellite instances.
func (l *TleLoader) Load(r io.Reader) ([]types.Satellite, error) {
	scanner := bufio.NewScanner(r)
	var satellites []types.Satellite

	for scanner.Scan() {
		line1 := strings.TrimSpace(scanner.Text())
		if line1 == "" {
			continue
		}

		var name string
		if !strings.HasPrefix(line1, "1") {
			name = line1
			if !scanner.Scan() {
				return nil, errors.New(errCannotParse)
			}
			line1 = strings.TrimSpace(scanner.Text())
			if !strings.HasPrefix(line1, "1") {
				return nil, errors.New(errCannotParse)
			}
		}

		if !scanner.Scan() {
			return nil, errors.New(errCannotParse)
		}
		line2 := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line2, "2") {
			return nil, errors.New(errCannotParse)
		}

		if name == "" && len(line1) > 5 {
			name = strings.TrimSpace(line1[2:6])
		}

		epochStr := strings.TrimSpace(line1[18:32])
		epoch, err := parseEpoch(epochStr)
		if err != nil {
			return nil, err
		}

		inclination, _ := strconv.ParseFloat(strings.TrimSpace(line2[8:16]), 64)
		raan, _ := strconv.ParseFloat(strings.TrimSpace(line2[17:25]), 64)
		eccentricity, _ := strconv.ParseFloat("0."+strings.TrimSpace(line2[26:33]), 64)
		argPerigee, _ := strconv.ParseFloat(strings.TrimSpace(line2[34:42]), 64)
		meanAnomaly, _ := strconv.ParseFloat(strings.TrimSpace(line2[43:51]), 64)
		meanMotion, _ := strconv.ParseFloat(strings.TrimSpace(line2[52:63]), 64)

		builder := l.satelliteBuilder
		builder.SetName(name).
			SetInclination(inclination).
			SetRightAscension(raan).
			SetEccentricity(eccentricity).
			SetArgumentOfPerigee(argPerigee).
			SetMeanAnomaly(meanAnomaly).
			SetMeanMotion(meanMotion).
			SetEpoch(epoch).
			ConfigureISL(func(b *links.IslProtocolBuilder) *links.IslProtocolBuilder {
				return b
			})

		sat := builder.Build()
		satellites = append(satellites, sat)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	log.Printf("Parsed %d satellites from TLE", len(satellites))
	return satellites, nil
}

// parseEpoch parses YYDDD.DDDDDDDD into a time.Time assuming year 2000+
func parseEpoch(epoch string) (time.Time, error) {
	if len(epoch) < 5 {
		return time.Time{}, fmt.Errorf("invalid epoch format: %s", epoch)
	}
	yearPrefix := 2000
	yy, err := strconv.Atoi(epoch[0:2])
	if err != nil {
		return time.Time{}, err
	}
	doy, err := strconv.ParseFloat(epoch[2:], 64)
	if err != nil {
		return time.Time{}, err
	}

	startOfYear := time.Date(yearPrefix+yy, 1, 1, 0, 0, 0, 0, time.UTC)
	return startOfYear.Add(time.Duration((doy - 1) * 24 * float64(time.Hour))), nil
}

// Ensure TleLoader implements SatelliteDataSourceLoader interface
var _ SatelliteDataSourceLoader = (*TleLoader)(nil)
