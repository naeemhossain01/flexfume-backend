package com.seamlance.perfume.service;

import com.seamlance.perfume.constants.ErrorConstant;
import com.seamlance.perfume.dao.DeliveryCostDao;
import com.seamlance.perfume.entity.DeliveryCost;
import com.seamlance.perfume.exception.ResourceNotFoundException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.dao.DataAccessException;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.List;

@Service
public class DeliveryCostServiceImpl implements DeliveryCostService {

    @Autowired
    private DeliveryCostDao deliveryCostDao;

    @Override
    public DeliveryCost addCost(DeliveryCost deliveryCost) {
        //TODO: validation needed;

        try {
            deliveryCost = deliveryCostDao.saveAndFlush(deliveryCost);
        } catch (DataAccessException e) {
            e.printStackTrace();
        }

        return deliveryCost;
    }

    @Override
    public DeliveryCost updateCost(String id, DeliveryCost deliveryCost) {
        DeliveryCost updatedDeliveryCost = deliveryCostDao.findById(id).orElseThrow(() -> new ResourceNotFoundException(ErrorConstant.DELIVERY_COST_INFO_NOT_FOUND));

        if(deliveryCost.getLocation() != null) updatedDeliveryCost.setLocation(deliveryCost.getLocation());
        if(deliveryCost.getService() != null) updatedDeliveryCost.setService(deliveryCost.getService());
        if(deliveryCost.getCost() != null) updatedDeliveryCost.setCost(deliveryCost.getCost());

        try {
            updatedDeliveryCost = deliveryCostDao.save(updatedDeliveryCost);
        } catch (DataAccessException e) {
            e.printStackTrace();
        }

        return updatedDeliveryCost;
    }

    @Override
    public List<DeliveryCost> getDeliveryCostByLocation(String location) {
        List<DeliveryCost> deliveryCostList = new ArrayList<>();

        try {
            deliveryCostList = deliveryCostDao.findByLocationContaining(location);
        } catch (DataAccessException e) {
            e.printStackTrace();
        }

        return deliveryCostList;
    }

    @Override
    public DeliveryCost getDeliveryCostById(String id) {
        return deliveryCostDao.findById(id).orElseThrow(() -> new ResourceNotFoundException(ErrorConstant.DELIVERY_COST_INFO_NOT_FOUND));
    }

    @Override
    public List<DeliveryCost> getAllDeliveryCost() {
        List<DeliveryCost> deliveryCostList = new ArrayList<>();

        try {
            deliveryCostList = deliveryCostDao.findAll();
        } catch (DataAccessException e) {
            e.printStackTrace();
        }

        return deliveryCostList;
    }

    @Override
    public void deleteDeliveryCost(String id) {
        DeliveryCost deliveryCost = this.getDeliveryCostById(id);

        deliveryCostDao.deleteById(id);
    }
}
