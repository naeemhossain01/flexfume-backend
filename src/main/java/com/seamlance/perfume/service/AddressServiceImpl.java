package com.seamlance.perfume.service;

import com.seamlance.perfume.constants.ErrorConstant;
import com.seamlance.perfume.dao.AddressDao;
import com.seamlance.perfume.entity.Address;
import com.seamlance.perfume.entity.User;
import com.seamlance.perfume.exception.ResourceNotFoundException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.dao.DataAccessException;
import org.springframework.stereotype.Service;

@Service
public class AddressServiceImpl implements AddressService {
    @Autowired
    private UserService userService;

    @Autowired
    private AddressDao addressDao;

    @Override
    public Address addAddress(Address address) {
        //TODO: Validation need;

        User user = userService.getLoginUser();

        if(user == null) {
            throw new ResourceNotFoundException(ErrorConstant.USER_NOT_FOUND);
        }

        Address existingAddress = user.getAddress();

        if(existingAddress == null) {
            address.setUser(user);
        }
        try {
            address = addressDao.saveAndFlush(address);
        } catch (DataAccessException e) {
            e.printStackTrace();
        }

        return address;
    }

    @Override
    public Address getAddressByUser(String userId) {
        return addressDao.findByUserId(userId).orElseThrow(() -> new ResourceNotFoundException(ErrorConstant.ADDRESS_NOT_FOUND));
    }

    @Override
    public Address updateAddress(String addressId, Address address) {
        User user = userService.getLoginUser();

        Address updatedAddress = addressDao.findByIdAndUserId(addressId, user.getId());

        if(updatedAddress == null) {
            throw new ResourceNotFoundException(ErrorConstant.ADDRESS_NOT_FOUND);
        }

        if(address.getBuildingName() != null) updatedAddress.setBuildingName(address.getBuildingName());
        if(address.getArea() != null) updatedAddress.setArea(address.getArea());
        if(address.getRoad() != null) updatedAddress.setRoad(address.getRoad());
        if(address.getCity() != null) updatedAddress.setCity(address.getCity());

        try {
            updatedAddress = addressDao.save(updatedAddress);
        } catch (DataAccessException e) {
            e.printStackTrace();
        }

        return updatedAddress;
    }
}
