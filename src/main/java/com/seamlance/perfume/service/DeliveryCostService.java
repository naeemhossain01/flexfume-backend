package com.seamlance.perfume.service;

import com.seamlance.perfume.entity.DeliveryCost;

import java.util.List;

public interface DeliveryCostService {
    DeliveryCost addCost(DeliveryCost deliveryCost);
    DeliveryCost updateCost(String id, DeliveryCost deliveryCost);
    List<DeliveryCost> getDeliveryCostByLocation(String location);
    DeliveryCost getDeliveryCostById(String id);
    List<DeliveryCost> getAllDeliveryCost();
    void deleteDeliveryCost(String id);
}
