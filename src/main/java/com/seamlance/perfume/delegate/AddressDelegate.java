package com.seamlance.perfume.delegate;

import com.seamlance.perfume.info.AddressInfo;

public interface AddressDelegate {
    AddressInfo addAddress(AddressInfo addressInfo);
    AddressInfo updateAddress(String userId, AddressInfo addressInfo);
    AddressInfo getAddressByUser(String userId);
}
