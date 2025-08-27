package com.seamlance.perfume.service;

import com.seamlance.perfume.entity.Address;

public interface AddressService {
    Address addAddress(Address address);

    Address getAddressByUser(String userId);

    Address updateAddress(String addressId, Address address);
}
