package com.seamlance.perfume.delegate;

import com.seamlance.perfume.info.DeliveryCostInfo;

import java.util.List;

public interface DeliveryCostDelegate {
    DeliveryCostInfo addCost(DeliveryCostInfo deliveryCostInfo);
    DeliveryCostInfo updateCost(String id, DeliveryCostInfo deliveryCostInfo);
    List<DeliveryCostInfo> getDeliveryCostByLocation(String location);
    List<DeliveryCostInfo> getAllDeliveryCost();
    DeliveryCostInfo getDeliveryCostById(String id);
    void deleteDeliveryCost(String id);
}
