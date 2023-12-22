package collector

import (
	"fmt"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"math/rand"
	"runtime"
	"sync"
)

func (c *Collector) collectAll() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.collectMem()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.collectCPU()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.collectMemStats()
	}()

	wg.Wait()
}

func (c *Collector) collectMemStats() {
	// Собираем метрики gauge
	runtime.ReadMemStats(c.memStats)
	c.muGauge.Lock()
	c.gaugeMetrics["Alloc"] = float64(c.memStats.Alloc)
	c.gaugeMetrics["BuckHashSys"] = float64(c.memStats.BuckHashSys)
	c.gaugeMetrics["Frees"] = float64(c.memStats.Frees)
	c.gaugeMetrics["GCCPUFraction"] = c.memStats.GCCPUFraction
	c.gaugeMetrics["GCSys"] = float64(c.memStats.GCSys)
	c.gaugeMetrics["HeapAlloc"] = float64(c.memStats.HeapAlloc)
	c.gaugeMetrics["HeapIdle"] = float64(c.memStats.HeapIdle)
	c.gaugeMetrics["HeapInuse"] = float64(c.memStats.HeapInuse)
	c.gaugeMetrics["HeapObjects"] = float64(c.memStats.HeapObjects)
	c.gaugeMetrics["HeapReleased"] = float64(c.memStats.HeapReleased)
	c.gaugeMetrics["HeapSys"] = float64(c.memStats.HeapSys)
	c.gaugeMetrics["LastGC"] = float64(c.memStats.LastGC)
	c.gaugeMetrics["Lookups"] = float64(c.memStats.Lookups)
	c.gaugeMetrics["MCacheInuse"] = float64(c.memStats.MCacheInuse)
	c.gaugeMetrics["MCacheSys"] = float64(c.memStats.MCacheSys)
	c.gaugeMetrics["MSpanInuse"] = float64(c.memStats.MSpanInuse)
	c.gaugeMetrics["MSpanSys"] = float64(c.memStats.MSpanSys)
	c.gaugeMetrics["Mallocs"] = float64(c.memStats.Mallocs)
	c.gaugeMetrics["NextGC"] = float64(c.memStats.NextGC)
	c.gaugeMetrics["NumForcedGC"] = float64(c.memStats.NumForcedGC)
	c.gaugeMetrics["NumGC"] = float64(c.memStats.NumGC)
	c.gaugeMetrics["OtherSys"] = float64(c.memStats.OtherSys)
	c.gaugeMetrics["PauseTotalNs"] = float64(c.memStats.PauseTotalNs)
	c.gaugeMetrics["StackInuse"] = float64(c.memStats.StackInuse)
	c.gaugeMetrics["StackSys"] = float64(c.memStats.StackSys)
	c.gaugeMetrics["Sys"] = float64(c.memStats.Sys)
	c.gaugeMetrics["TotalAlloc"] = float64(c.memStats.TotalAlloc)
	c.gaugeMetrics["RandomValue"] = rand.Float64()
	c.muGauge.Unlock()

	// Собираем метики counter
	c.muCounter.Lock()
	c.pollCount += 1
	c.muCounter.Unlock()
}

func (c *Collector) collectMem() {
	vm, err := mem.VirtualMemory()
	if err != nil {
		logger.Error().Err(fmt.Errorf("error due collectPS VM %w", err))
		return
	}
	c.muGauge.Lock()
	c.gaugeMetrics["TotalMemory"] = float64(vm.Total)
	c.gaugeMetrics["FreeMemory"] = float64(vm.Free)
	c.muGauge.Unlock()
}

func (c *Collector) collectCPU() {
	cp, err := cpu.Percent(0, true)
	if err != nil {
		logger.Error().Err(fmt.Errorf("error due cpu usage %w", err))
		return
	}
	c.muGauge.Lock()
	for i, utl := range cp {
		c.gaugeMetrics[fmt.Sprintf("CPUutilization%d", i)] = utl
	}
	c.muGauge.Unlock()
}
