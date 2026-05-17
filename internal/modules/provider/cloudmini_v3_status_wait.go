package provider

import (
	"context"
	"strings"
	"time"
)

func (adapter *CloudminiV3Adapter) waitForUsableProxy(ctx context.Context, operation OperationContext, client *cloudminiV3Client, initial cloudminiV3Proxy, resourceID string) (cloudminiV3Proxy, error) {
	if cloudminiV3ProxyStatusUsable(initial.Status) {
		return initial, nil
	}
	pollCtx, cancel := context.WithTimeout(ctx, adapter.pollTimeout)
	defer cancel()
	ticker := time.NewTicker(adapter.pollInterval)
	defer ticker.Stop()
	lastProxy := initial
	for {
		readID := strings.TrimSpace(lastProxy.ID)
		if readID == "" {
			readID = strings.TrimSpace(resourceID)
		}
		if readID == "" {
			adapterErr, _ := cloudminiV3ProxyStatusNotUsable(lastProxy.Status)
			return lastProxy, adapterErr
		}
		proxy, err := client.getProxy(pollCtx, operation, readID)
		if err != nil {
			if pollCtx.Err() != nil {
				adapterErr, _ := cloudminiV3ProxyStatusNotUsable(lastProxy.Status)
				return lastProxy, adapterErr
			}
			return lastProxy, err
		}
		lastProxy = proxy
		if cloudminiV3ProxyStatusUsable(proxy.Status) {
			return proxy, nil
		}
		select {
		case <-pollCtx.Done():
			adapterErr, _ := cloudminiV3ProxyStatusNotUsable(lastProxy.Status)
			return lastProxy, adapterErr
		case <-ticker.C:
		}
	}
}
