import React from 'react';
import PropTypes from 'prop-types';
import styled from 'styled-components';
import { rem } from 'polished';

import logo from './logo-mdl.png';
import logoWhite from './logo-white.png';


const Img = styled.img.attrs({
  alt: 'MDL',
})`
  height: ${rem(40)};
  max-width: 100%;
`;

const Logo = props => (
  <a href="https://mdl.life">
    <Img {...props} src={props.white ? logoWhite :  logo} />
  </a>
);

Logo.propTypes = {
  white: PropTypes.bool,
};

Logo.defaultProps = {
  white: false,
};

export default Logo;
