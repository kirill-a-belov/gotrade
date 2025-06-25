package chart

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/kirill-a-belov/trader/internal/chart/model"
	"github.com/kirill-a-belov/trader/pkg/tracer"
)

func (m *Module) Serve(ctx context.Context) error {
	_, span := tracer.Start(ctx, "pkg.internal.chart.Serve")
	defer span.End()

	const htmlDirPath = "./web/static"

	http.Handle("/", http.FileServer(http.Dir(htmlDirPath)))
	http.HandleFunc("/candles", m.candlesHandler)
	http.HandleFunc("/markers", m.markersHandler)

	var err error
	go func() {
		if err = http.ListenAndServe(fmt.Sprintf(":%s", m.config.Port), nil); err != nil {
			err = errors.Wrap(err, "ListenAndServe")
		}
	}()

	return err
}

func (m *Module) candlesHandler(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "pkg.internal.chart.candlesHandler")
	defer span.End()

	if err := json.NewEncoder(w).Encode(m.candleList); err != nil {
		m.log.Error(r.Context(), fmt.Sprintf("error in candles handler: %v", err))
	}
}

func (m *Module) markersHandler(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "pkg.internal.chart.markersHandler")
	defer span.End()

	wholeMarkerList := make([]*model.Marker, 0)
	for _, markerList := range m.markerStorage {
		wholeMarkerList = append(wholeMarkerList, markerList...)
	}

	if err := json.NewEncoder(w).Encode(wholeMarkerList); err != nil {
		m.log.Error(r.Context(), fmt.Sprintf("error in markers handler: %v", err))
	}
}
