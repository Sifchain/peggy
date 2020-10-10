pragma solidity ^0.5.0;

/**
 * @title WhiteList
 * @dev WhiteList contract records the ERC 20 list that can be locked in BridgeBank.
 **/

contract WhiteList {
    mapping(address => bool) whiteList;

    /*
     * @dev: Event declarations
     */
    event LogWhiteListUpdate(address _token, bool _value);

    constructor() public {
        whiteList[address(0)] = true;
    }

    /*
     * @dev: Set the token address in whitelist
     *
     * @param _token: ERC 20's address
     * @param _inList: set the _token in list or not
     * @return: new value of if _token in whitelist
     */
    function setTokenInWhiteList(address _token, bool _inList)
        internal
        returns (bool)
    {
        whiteList[_token] = _inList;
        emit LogWhiteListUpdate(_token, _inList);
        return _inList;
    }

    /*
     * @dev: Get if the token in whitelist
     *
     * @param _token: ERC 20's address
     * @return: if _token in whitelist
     */
    function getTokenInWhiteList(address _token) public view returns (bool) {
        return whiteList[_token];
    }
}