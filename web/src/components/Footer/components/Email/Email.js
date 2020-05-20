import React from 'react';
import styled from 'styled-components';

const CONTACT_EMAIL_1 = 'rt@mdl.life';
const CONTACT_EMAIL_2 = 'tm@mdl.life';

const Email = styled.a`
  color: inherit;
  text-decoration: none;

  &:hover {
    text-decoration: underline;
  }
`;

export default () => (
  <div>
    <Email href={`mailto:${CONTACT_EMAIL_1}`}>
      {CONTACT_EMAIL_1}
    </Email>
    <br/>
    <Email href={`mailto:${CONTACT_EMAIL_2}`}>
     {CONTACT_EMAIL_2}
    </Email>
  </div>
);
