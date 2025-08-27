package com.seamlance.perfume.delegate;

import com.seamlance.perfume.info.DeliveryCostInfo;
import com.seamlance.perfume.mapper.DeliveryCostMapper;
import com.seamlance.perfume.service.DeliveryCostService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import java.util.List;

@Component
public class DeliveryCostDelegateImpl implements DeliveryCostDelegate {

    @Autowired
    private DeliveryCostService deliveryCostService;

    @Autowired
    private DeliveryCostMapper deliveryCostMapper;

    @Override
    public DeliveryCostInfo addCost(DeliveryCostInfo deliveryCostInfo) {
        return deliveryCostMapper.toInfo(deliveryCostService.addCost(deliveryCostMapper.toEntity(deliveryCostInfo)));
    }

    @Override
    public DeliveryCostInfo updateCost(String id, DeliveryCostInfo deliveryCostInfo) {
        return deliveryCostMapper.toInfo(deliveryCostService.updateCost(id, deliveryCostMapper.toEntity(deliveryCostInfo)));
    }

    @Override
    public List<DeliveryCostInfo> getDeliveryCostByLocation(String location) {
        return deliveryCostMapper.toInfoList(deliveryCostService.getDeliveryCostByLocation(location));
    }

    @Override
    public List<DeliveryCostInfo> getAllDeliveryCost() {
        return deliveryCostMapper.toInfoList(deliveryCostService.getAllDeliveryCost());
    }

    @Override
    public DeliveryCostInfo getDeliveryCostById(String id) {
        return deliveryCostMapper.toInfo(deliveryCostService.getDeliveryCostById(id));
    }

    @Override
    public void deleteDeliveryCost(String id) {
        deliveryCostService.deleteDeliveryCost(id);
    }
}
