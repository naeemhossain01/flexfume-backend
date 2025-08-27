package com.seamlance.perfume.delegate;

import com.seamlance.perfume.info.AddressInfo;
import com.seamlance.perfume.mapper.AddressMapper;
import com.seamlance.perfume.service.AddressService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

@Component
public class AddressDelegateImpl implements AddressDelegate {
    @Autowired
    private AddressMapper addressMapper;

    @Autowired
    private AddressService addressService;


    @Override
    public AddressInfo addAddress(AddressInfo addressInfo) {
        return addressMapper.toInfo(addressService.addAddress(addressMapper.toEntity(addressInfo)));
    }

    @Override
    public AddressInfo updateAddress(String userId, AddressInfo addressInfo) {
        return addressMapper.toInfo(addressService.updateAddress(userId, addressMapper.toEntity(addressInfo)));
    }

    @Override
    public AddressInfo getAddressByUser(String userId) {
        return addressMapper.toInfo(addressService.getAddressByUser(userId));
    }
}
