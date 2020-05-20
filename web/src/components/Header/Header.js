import React from 'react';
import styled from 'styled-components';
import PropTypes from 'prop-types';
import { rem } from 'polished';

import { SPACE, COLORS } from 'config';
import Container from '../Container';
import Logo from '../Logo';

const Wrapper = styled.div`
  padding: ${rem(SPACE[6])} 0;
  width: 100%;
  border-bottom: ${props => (props.border ? `2px solid ${COLORS.gray[1]}` : 'none')}
`;

const Centered = styled.div`
  text-align: center;
  align-content: center;
  width: 100%;
`

const Header = ({ white, border }) => (
  <Wrapper border={border}>
    <Container>
      <Centered>
        <Logo white={white} />
      </Centered>
    </Container>
  </Wrapper>
);

Header.propTypes = {
  white: PropTypes.bool,
  border: PropTypes.bool,
};

Header.defaultProps = {
  white: false,
  border: false,
};

export default Header;
